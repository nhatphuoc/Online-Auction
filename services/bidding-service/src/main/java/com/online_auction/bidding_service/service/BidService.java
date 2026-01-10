package com.online_auction.bidding_service.service;

import lombok.RequiredArgsConstructor;

import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageImpl;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import com.online_auction.bidding_service.client.ProductServiceClient;
import com.online_auction.bidding_service.config.RabbitMQConfig;
import com.online_auction.bidding_service.domain.AutoBid;
import com.online_auction.bidding_service.domain.BiddingHistory;
import com.online_auction.bidding_service.domain.BiddingHistory.BidStatus;
import com.online_auction.bidding_service.domain.Product;
import com.online_auction.bidding_service.dto.request.ProductBidRequest;
import com.online_auction.bidding_service.dto.response.ApiResponse;
import com.online_auction.bidding_service.dto.response.BiddingHistorySearchResponse;
import com.online_auction.bidding_service.dto.response.ProductBidSuccessData;
import com.online_auction.bidding_service.dto.response.UserBidResponse;
import com.online_auction.bidding_service.event.BidPlacedEvent;
import com.online_auction.bidding_service.repository.AutoBidRepository;
import com.online_auction.bidding_service.repository.BiddingHistoryRepository;
import com.online_auction.bidding_service.repository.ProductRepository;
import com.online_auction.bidding_service.specs.BiddingHistorySpecs;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class BidService {

        private final BiddingHistoryRepository biddingHistoryRepository;
        private final ProductRepository productRepository;
        private final ProductServiceClient productServiceClient;
        private final RabbitTemplate rabbitTemplate;
        private final AutoBidRepository autoBidRepository;

        @Transactional
        public ApiResponse<?> placeBid(
                        Long productId,
                        Long bidderId,
                        Double amount,
                        String requestId,
                        String userJwt) {

                ProductBidRequest req = new ProductBidRequest(bidderId, amount, requestId);

                // 1. Call Product-Service atomic endpoint
                ApiResponse<ProductBidSuccessData> resp = productServiceClient.placeBidToProductService(productId, req,
                                userJwt);

                System.out.println("Response: " + resp);
                if (resp == null) {
                        saveHistory(productId, bidderId, amount, requestId, BidStatus.FAILED,
                                        "NULL_RESPONSE_FROM_PRODUCT_SERVICE");
                        return ApiResponse.fail("Product service returned empty response");
                }

                // 2. Nếu lỗi từ product-service
                if (!resp.isSuccess()) {
                        saveHistory(productId, bidderId, amount, requestId, BidStatus.FAILED, resp.getMessage());
                        return ApiResponse.fail(resp.getMessage());
                }

                // 3. Thành công
                ProductBidSuccessData data = resp.getData();
                saveHistory(productId, bidderId, amount, requestId, BidStatus.SUCCESS, null);

                // 4. Publish event
                BidPlacedEvent event = new BidPlacedEvent(
                                productId,
                                bidderId,
                                amount,
                                data != null ? data.getPreviousHighestBidder() : null,
                                requestId);

                rabbitTemplate.convertAndSend(
                                RabbitMQConfig.EXCHANGE,
                                RabbitMQConfig.ROUTING_KEY_BID_SUCCESS,
                                event);
                Product product = productRepository
                                .findByIdForUpdate(productId)
                                .orElseThrow(() -> new RuntimeException("PRODUCT_NOT_FOUND"));
                triggerAutoBid(product, bidderId);

                return ApiResponse.ok(data, resp.getMessage());
        }

        @Transactional
        private void triggerAutoBid(Product product, Long lastBidderId) {

                List<AutoBid> autoBids = autoBidRepository
                                .findByProductIdAndActiveTrueOrderByMaxAmountDesc(product.getId());

                if (autoBids.isEmpty())
                        return;

                for (AutoBid autoBid : autoBids) {

                        // 1. Không auto bid chính mình
                        if (autoBid.getBidderId().equals(lastBidderId))
                                continue;

                        double currentPrice = product.getCurrentPrice() != null
                                        ? product.getCurrentPrice()
                                        : product.getStartingPrice();

                        double nextPrice = currentPrice + product.getStepPrice();

                        // 2. Nếu vượt max → deactivate
                        if (nextPrice > autoBid.getMaxAmount()) {
                                autoBid.setActive(false);
                                autoBidRepository.save(autoBid);
                                continue;
                        }

                        // 3. UPDATE PRODUCT LOCAL (CỐT LÕI)
                        Long prevBidder = product.getCurrentBidder();

                        product.setCurrentPrice(nextPrice);
                        product.setCurrentBidder(autoBid.getBidderId());
                        product.setBidCount(product.getBidCount() + 1);

                        productRepository.save(product);

                        // 4. SAVE HISTORY (AUTO)
                        saveHistory(
                                        product.getId(),
                                        autoBid.getBidderId(),
                                        nextPrice,
                                        "AUTO_" + UUID.randomUUID(),
                                        BidStatus.SUCCESS,
                                        null);

                        // 5. Chỉ 1 auto bid mỗi vòng
                        break;
                }
        }

        private void saveHistory(Long productId, Long bidderId, Double amount,
                        String requestId, BidStatus status, String reason) {
                BiddingHistory h = new BiddingHistory();
                h.setProductId(productId);
                h.setBidderId(bidderId);
                h.setAmount(amount);
                h.setRequestId(requestId);
                h.setStatus(status);
                h.setReason(reason);
                h.setCreatedAt(LocalDateTime.now());
                biddingHistoryRepository.save(h);
        }

        public Page<BiddingHistorySearchResponse> search(
                        Long productId,
                        Long bidderId,
                        BiddingHistory.BidStatus status,
                        String requestId,
                        LocalDateTime from,
                        LocalDateTime to,
                        Pageable pageable) {
                return biddingHistoryRepository
                                .findAll(BiddingHistorySpecs.search(productId, bidderId, status, requestId, from, to),
                                                pageable)
                                .map(a -> toResponse(null, requestId));
        }

        private BiddingHistorySearchResponse toResponse(BiddingHistory entity, String bidderName) {
                return new BiddingHistorySearchResponse(
                                entity.getId(),
                                entity.getProductId(),
                                entity.getBidderId(),
                                bidderName,
                                entity.getAmount(),
                                entity.getRequestId(),
                                entity.getStatus(),
                                entity.getReason(),
                                entity.getCreatedAt());
        }

        public List<UserBidResponse> getUserBids(Long userId) {
                List<BiddingHistory> bids = biddingHistoryRepository.findAllByBidderId(userId);

                return bids.stream().map(bid -> {
                        Product product = productRepository.findById(bid.getProductId())
                                        .orElseThrow(() -> new RuntimeException("Product not found"));

                        return UserBidResponse.builder()
                                        .bidId(bid.getId())
                                        .bidAmount(bid.getAmount())
                                        .bidStatus(bid.getStatus().name())
                                        .reason(bid.getReason())
                                        .bidCreatedAt(bid.getCreatedAt())
                                        .productId(product.getId())
                                        .productName(product.getName())
                                        .thumbnailUrl(product.getThumbnailUrl())
                                        .currentPrice(product.getCurrentPrice())
                                        .buyNowPrice(product.getBuyNowPrice())
                                        .endAt(product.getEndAt())
                                        .autoExtend(product.isAutoExtend())
                                        .currentBidder(product.getCurrentBidder())
                                        .build();
                }).toList();
        }

        public Page<UserBidResponse> getUserBidsFiltered(Long userId, String filter, int page, int size) {
                List<UserBidResponse> allBids = getUserBids(userId);

                List<UserBidResponse> filteredBids = switch (filter.toLowerCase()) {
                        case "winning" -> allBids.stream()
                                        .filter(bid -> bid.getCurrentPrice() != null && bid.getBidAmount() != null
                                                        && bid.getBidAmount().equals(bid.getCurrentPrice()))
                                        .toList();
                        case "outbid" -> allBids.stream()
                                        .filter(bid -> bid.getCurrentPrice() != null && bid.getBidAmount() != null
                                                        && !bid.getBidAmount().equals(bid.getCurrentPrice()))
                                        .toList();
                        default -> allBids;
                };

                // Tạo Page từ List
                int start = Math.min(page * size, filteredBids.size());
                int end = Math.min(start + size, filteredBids.size());
                List<UserBidResponse> pageContent = filteredBids.subList(start, end);

                return new PageImpl<>(pageContent, PageRequest.of(page, size), filteredBids.size());
        }

        public Page<BiddingHistorySearchResponse> getBidsByProduct(
                        Long productId,
                        Long bidderId,
                        BiddingHistory.BidStatus status,
                        int page,
                        int size) {
                Pageable pageable = PageRequest.of(
                                page,
                                size,
                                Sort.by(Sort.Direction.DESC, "createdAt"));

                return biddingHistoryRepository.searchByProductId(
                                productId,
                                pageable);
        }

}
