package com.online_auction.bidding_service.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageImpl;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.online_auction.bidding_service.client.ProductServiceClient;
import com.online_auction.bidding_service.config.RabbitMQConfig;
import com.online_auction.bidding_service.domain.AutoBid;
import com.online_auction.bidding_service.domain.BiddingHistory;
import com.online_auction.bidding_service.domain.BiddingHistory.BidStatus;
import com.online_auction.bidding_service.domain.Product;
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
@Slf4j
public class BidService {
        private ObjectMapper objectMapper = new ObjectMapper();
        private final BiddingHistoryRepository biddingHistoryRepository;
        private final ProductRepository productRepository;
        private final ProductServiceClient productServiceClient;
        private final RabbitTemplate rabbitTemplate;
        private final AutoBidRepository autoBidRepository;

        @Transactional
        public ApiResponse<?> placeBid(
                        Long productId,
                        Long bidderId,
                        Double bidAmount,
                        String requestId) {

                // ====== 1. Lock product ======
                Product product = productRepository.findByIdForUpdate(productId)
                                .orElse(null);

                if (product == null) {
                        saveHistory(productId, bidderId, bidAmount, requestId,
                                        BidStatus.FAILED, "PRODUCT_NOT_FOUND");
                        return ApiResponse.fail("Product not found");
                }

                // ====== 2. Check auction end ======
                if (LocalDateTime.now().isAfter(product.getEndAt())) {
                        saveHistory(productId, bidderId, bidAmount, requestId,
                                        BidStatus.FAILED, "AUCTION_ENDED");
                        return ApiResponse.fail("Auction has already ended");
                }

                Long previousHighestBidder = product.getCurrentBidder();

                // ====== 3. Case: Chưa có ai bid ======
                if (product.getCurrentPrice() == null) {

                        if (bidAmount < product.getStartingPrice()) {
                                saveHistory(productId, bidderId, bidAmount, requestId,
                                                BidStatus.FAILED, "LOWER_THAN_STARTING_PRICE");
                                return ApiResponse.fail("Bid amount must be equal or higher than starting price");
                        }

                        product.setCurrentPrice(bidAmount);
                        product.setCurrentBidder(bidderId);
                        product.setBidCount(product.getBidCount() + 1);

                        productRepository.save(product);

                        saveHistory(productId, bidderId, bidAmount, requestId,
                                        BidStatus.SUCCESS, null);

                        publishBidSuccessEvent(
                                        productId,
                                        bidderId,
                                        bidAmount,
                                        null,
                                        requestId);

                        return ApiResponse.ok(
                                        new ProductBidSuccessData(bidAmount, null),
                                        "Bid placed successfully");
                }

                // ====== 4. Case: Đã có người bid ======
                if (bidAmount <= product.getCurrentPrice()) {
                        return ApiResponse.fail("Bid amount must be higher than current price");
                }

                List<AutoBid> autoBids = autoBidRepository
                                .findByProductIdAndActiveTrueOrderByMaxAmountDesc(productId);
                AutoBid highest = autoBids.stream()
                                .filter(ab -> !ab.getBidderId().equals(bidderId))
                                .findFirst()
                                .orElse(null);
                Long winningBidderId = bidderId;
                if (highest != null) {
                        if (bidAmount < highest.getMaxAmount()) {
                                bidAmount = bidAmount + product.getStepPrice();
                                winningBidderId = highest.getBidderId();
                        }
                }

                // ====== 5. Update product ======
                product.setCurrentPrice(bidAmount);
                product.setCurrentBidder(winningBidderId);
                product.setBidCount(product.getBidCount() + 1);

                productRepository.save(product);

                // ====== 6. Save history SUCCESS ======
                saveHistory(productId, winningBidderId, bidAmount, requestId,
                                BidStatus.SUCCESS, null);

                // ====== 7. Publish event ======
                publishBidSuccessEvent(
                                productId,
                                winningBidderId,
                                bidAmount,
                                previousHighestBidder,
                                requestId);

                return ApiResponse.ok(
                                new ProductBidSuccessData(bidAmount, previousHighestBidder),
                                "Bid placed successfully");
        }

        public void publishBidSuccessEvent(
                        Long productId,
                        Long bidderId,
                        Double amount,
                        Long previousHighestBidder,
                        String requestId) {
                BidPlacedEvent event = new BidPlacedEvent(
                                productId,
                                bidderId,
                                amount,
                                previousHighestBidder,
                                requestId);

                rabbitTemplate.convertAndSend(
                                RabbitMQConfig.EXCHANGE,
                                RabbitMQConfig.ROUTING_KEY_BID_SUCCESS,
                                event);
        }

        @Transactional
        private void triggerAutoBid(Product product, Long lastBidderId) {

                List<AutoBid> autoBids = autoBidRepository
                                .findByProductIdAndActiveTrueOrderByMaxAmountDesc(product.getId());
                System.out.println("size: " + autoBids.size());
                if (autoBids.isEmpty()) {
                        return;
                }

                AutoBid highest = autoBids.get(0);
                System.out.println("HighestId: " + highest.getBidderId());
                System.out.println("currentId: " + lastBidderId);

                System.out.println("HighestAmount: " + highest.getMaxAmount());
                System.out.println("CurrentAmount: " + product.getCurrentPrice());
                // Không auto-bid chính người vừa bid
                if (highest.getBidderId().equals(lastBidderId) && highest.getMaxAmount() != product.getCurrentPrice()) {
                        return;
                }

                double step = product.getStepPrice();
                double nextPrice = product.getCurrentPrice() + step;
                System.out.println("next price: " + nextPrice);

                // Nếu vượt quá maxAmount của autoBid → stop
                if (nextPrice > highest.getMaxAmount()) {
                        return;
                }

                Long prevBidder = product.getCurrentBidder();

                product.setCurrentPrice(nextPrice);
                product.setCurrentBidder(highest.getBidderId());
                product.setBidCount(product.getBidCount() + 1);

                productRepository.save(product);

                String bidId = "auto_bid_" + UUID.randomUUID();

                saveHistory(
                                product.getId(),
                                highest.getBidderId(),
                                nextPrice,
                                bidId,
                                BidStatus.SUCCESS,
                                null);

                publishBidSuccessEvent(
                                product.getId(),
                                highest.getBidderId(),
                                nextPrice,
                                prevBidder,
                                bidId);
        }

        public void saveHistory(Long productId, Long bidderId, Double amount,
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

        @Transactional
        public ApiResponse<?> cancelTopBid(Long productId, Long sellerId) {
                // ====== 1. Lock product ======
                Product product = productRepository.findByIdForUpdate(productId)
                                .orElse(null);

                if (product == null) {
                        return ApiResponse.fail("Product not found");
                }

                // ====== 2. Check seller authorization ======
                if (!product.getSellerId().equals(sellerId)) {
                        return ApiResponse.fail("You are not authorized to cancel bids for this product");
                }

                // ====== 3. Check if there is a current bidder ======
                if (product.getCurrentBidder() == null) {
                        return ApiResponse.fail("No bids to cancel");
                }

                // ====== 4. Find the latest SUCCESS bid for this product ======
                Pageable pageable = PageRequest.of(0, 1, Sort.by(Sort.Direction.DESC, "createdAt"));
                Page<BiddingHistory> latestSuccessBidPage = biddingHistoryRepository.findByProductIdAndStatus(
                                productId,
                                BidStatus.SUCCESS,
                                pageable);

                if (latestSuccessBidPage.isEmpty()) {
                        return ApiResponse.fail("No successful bid found to cancel");
                }

                BiddingHistory latestBid = latestSuccessBidPage.getContent().get(0);
                latestBid.setStatus(BidStatus.FAILED);
                latestBid.setReason("CANCELLED_BY_SELLER");
                biddingHistoryRepository.save(latestBid);

                // ====== 5. Find the previous highest bid (after the cancelled one) ======
                Pageable prevPageable = PageRequest.of(0, 1, Sort.by(Sort.Direction.DESC, "createdAt"));
                Page<BiddingHistory> prevBidPage = biddingHistoryRepository.findByProductIdAndStatus(
                                productId,
                                BidStatus.SUCCESS,
                                prevPageable);

                if (prevBidPage.hasContent()) {
                        BiddingHistory previousBid = prevBidPage.getContent().get(0);
                        product.setCurrentPrice(previousBid.getAmount());
                        product.setCurrentBidder(previousBid.getBidderId());
                } else {
                        // No previous bid, revert to starting price
                        product.setCurrentPrice(null);
                        product.setCurrentBidder(null);
                        product.setBidCount(0L);
                }

                // ====== 6. Save to products table in bidding-service database ======
                productRepository.save(product);

                // ====== 7. Publish event (optional, for notifications) ======
                // You can add a RabbitMQ event here if needed

                return ApiResponse.ok(null, "Top bid cancelled successfully");
        }

}
