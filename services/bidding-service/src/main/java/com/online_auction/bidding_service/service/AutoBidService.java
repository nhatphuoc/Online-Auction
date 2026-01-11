package com.online_auction.bidding_service.service;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Propagation;

import com.online_auction.bidding_service.domain.AutoBid;
import com.online_auction.bidding_service.domain.BiddingHistory.BidStatus;
import com.online_auction.bidding_service.domain.Product;
import com.online_auction.bidding_service.repository.AutoBidRepository;
import com.online_auction.bidding_service.repository.ProductRepository;

import org.springframework.transaction.annotation.Transactional;

@Service
public class AutoBidService {

        @Autowired
        private AutoBidRepository autoBidRepository;

        @Autowired
        private ProductRepository productRepository;

        @Autowired
        private BidService bidService;

        @Transactional

        public void registerAutoBid(Long productId, Long bidderId, Double maxAmount) {

                Product product = productRepository.findById(productId)
                                .orElseThrow(() -> new RuntimeException("PRODUCT_NOT_FOUND"));

                if (LocalDateTime.now().isAfter(product.getEndAt())) {
                        throw new IllegalStateException("AUCTION_ENDED");
                }

                double currentPrice = product.getCurrentPrice() != null
                                ? product.getCurrentPrice()
                                : product.getStartingPrice();

                if (maxAmount <= currentPrice) {
                        throw new IllegalArgumentException("MAX_AMOUNT_TOO_LOW");
                }

                // 1️⃣ Save / update auto-bid hiện tại
                AutoBid newAutoBid = autoBidRepository
                                .findByProductIdAndBidderId(productId, bidderId)
                                .orElse(
                                                AutoBid.builder()
                                                                .productId(productId)
                                                                .bidderId(bidderId)
                                                                .createdAt(LocalDateTime.now())
                                                                .build());

                newAutoBid.setMaxAmount(maxAmount);
                newAutoBid.setActive(true);
                newAutoBid.setUpdatedAt(LocalDateTime.now());

                autoBidRepository.save(newAutoBid);

                // 2️⃣ Lấy highest auto-bid KHÁC mình
                List<AutoBid> autoBids = autoBidRepository
                                .findByProductIdAndActiveTrueOrderByMaxAmountDesc(productId);

                AutoBid highest = autoBids.stream()
                                .filter(ab -> !ab.getBidderId().equals(bidderId))
                                .findFirst()
                                .orElse(null);

                double step = product.getStepPrice();
                double bidPrice;

                // Nếu chưa có đối thủ → không cần auto-bid
                if (highest == null) {
                        bidPrice = product.getStartingPrice();
                        bidService.placeBid(productId, bidderId, bidPrice, "auto_bid_" + UUID.randomUUID());
                        return;
                }

                // 3️⃣ CORE LOGIC (đúng như bạn mô tả)
                if (newAutoBid.getMaxAmount() < highest.getMaxAmount()) {
                        bidPrice = newAutoBid.getMaxAmount() + step;
                        bidService.placeBid(
                                        productId,
                                        highest.getBidderId(),
                                        bidPrice,
                                        "auto_bid_" + UUID.randomUUID());
                } else {
                        bidPrice = highest.getMaxAmount() + step;
                        bidService.placeBid(
                                        productId,
                                        newAutoBid.getBidderId(),
                                        bidPrice,
                                        "auto_bid_" + UUID.randomUUID());
                }
        }

}
