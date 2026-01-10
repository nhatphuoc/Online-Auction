package com.online_auction.bidding_service.service;

import java.time.LocalDateTime;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.online_auction.bidding_service.domain.AutoBid;
import com.online_auction.bidding_service.domain.Product;
import com.online_auction.bidding_service.repository.AutoBidRepository;
import com.online_auction.bidding_service.repository.ProductRepository;

import jakarta.transaction.Transactional;

@Service
public class AutoBidService {

    @Autowired
    private AutoBidRepository autoBidRepository;

    @Autowired
    private ProductRepository productRepository;

    @Transactional
    public void registerAutoBid(Long productId, Long bidderId, Double maxAmount) {

        Product product = productRepository.findById(productId)
                .orElseThrow(() -> new RuntimeException("PRODUCT_NOT_FOUND"));

        if (LocalDateTime.now().isAfter(product.getEndAt())) {
            throw new IllegalStateException("AUCTION_ENDED");
        }

        double minAllowed = product.getCurrentPrice() != null
                ? product.getCurrentPrice()
                : product.getStartingPrice();

        if (maxAmount <= minAllowed) {
            throw new IllegalArgumentException("MAX_AMOUNT_TOO_LOW");
        }

        AutoBid autoBid = autoBidRepository
                .findByProductIdAndBidderId(productId, bidderId)
                .orElse(
                        AutoBid.builder()
                                .productId(productId)
                                .bidderId(bidderId)
                                .createdAt(LocalDateTime.now())
                                .build());

        autoBid.setMaxAmount(maxAmount);
        autoBid.setActive(true);
        autoBid.setUpdatedAt(LocalDateTime.now());

        autoBidRepository.save(autoBid);
    }

}
