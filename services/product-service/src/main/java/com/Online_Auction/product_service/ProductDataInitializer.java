package com.Online_Auction.product_service;

import java.time.LocalDateTime;
import java.util.List;

import org.springframework.boot.CommandLineRunner;
import org.springframework.stereotype.Component;

import com.Online_Auction.product_service.domain.Product;
import com.Online_Auction.product_service.repository.ProductRepository;

import lombok.RequiredArgsConstructor;

@Component
@RequiredArgsConstructor
public class ProductDataInitializer implements CommandLineRunner {

        private final ProductRepository productRepository;

        @Override
        public void run(String... args) {

                if (productRepository.count() > 0) {
                        return; // prevent duplicate inserts
                }

                LocalDateTime now = LocalDateTime.now();

                List<Product> products = List.of(

                                // ================= PRODUCT 1 (EXPIRED, HAS WINNER) =================
                                Product.builder()
                                                .sellerId(10L)
                                                .name("Vintage Mechanical Watch")
                                                .thumbnailUrl("https://example.com/thumb/watch.jpg")
                                                .description("Classic vintage mechanical watch in good condition.")
                                                .parentCategoryId(1L)
                                                .parentCategoryName("Accessories")
                                                .categoryId(4L)
                                                .categoryName("Watches")
                                                .startingPrice(100.0)
                                                .currentPrice(180.0)
                                                .buyNowPrice(300.0)
                                                .stepPrice(10.0)
                                                .createdAt(now.minusDays(3))
                                                .endAt(now.minusHours(1))
                                                .autoExtend(false)
                                                .bidCount(Long.valueOf(5))
                                                .currentBidder(7L)
                                                .orderCreated(false)
                                                .sentEmail(false)
                                                .build(),

                                // ================= PRODUCT 2 =================
                                Product.builder()
                                                .sellerId(11L)
                                                .name("Gaming Mechanical Keyboard")
                                                .thumbnailUrl("https://example.com/thumb/keyboard.jpg")
                                                .description("RGB mechanical keyboard, lightly used.")
                                                .parentCategoryId(2L)
                                                .parentCategoryName("Electronics")
                                                .categoryId(7L)
                                                .categoryName("Computer Accessories")
                                                .startingPrice(50.0)
                                                .currentPrice(95.0)
                                                .buyNowPrice(150.0)
                                                .stepPrice(5.0)
                                                .createdAt(now.minusDays(2))
                                                .endAt(now.minusMinutes(30))
                                                .autoExtend(true)
                                                .bidCount(Long.valueOf(8))
                                                .currentBidder(8L)
                                                .orderCreated(false)
                                                .sentEmail(false)
                                                .build(),

                                // ================= PRODUCT 3 (NO BID) =================
                                Product.builder()
                                                .sellerId(10L)
                                                .name("Handmade Leather Wallet")
                                                .thumbnailUrl("https://example.com/thumb/wallet.jpg")
                                                .description("Handcrafted leather wallet, brand new.")
                                                .parentCategoryId(1L)
                                                .parentCategoryName("Fashion")
                                                .categoryId(5L)
                                                .categoryName("Wallets")
                                                .startingPrice(40.0)
                                                .currentPrice(null)
                                                .buyNowPrice(80.0)
                                                .stepPrice(5.0)
                                                .createdAt(now.minusDays(1))
                                                .endAt(now.minusMinutes(10))
                                                .autoExtend(false)
                                                .bidCount(Long.valueOf(0))
                                                .currentBidder(null)
                                                .orderCreated(false)
                                                .sentEmail(false)
                                                .build(),

                                // ================= PRODUCT 4 (NOT EXPIRED) =================
                                Product.builder()
                                                .sellerId(11L)
                                                .name("Wireless Noise Cancelling Headphones")
                                                .thumbnailUrl("https://example.com/thumb/headphones.jpg")
                                                .description("Premium wireless headphones with ANC.")
                                                .parentCategoryId(2L)
                                                .parentCategoryName("Electronics")
                                                .categoryId(8L)
                                                .categoryName("Audio")
                                                .startingPrice(200.0)
                                                .currentPrice(230.0)
                                                .buyNowPrice(350.0)
                                                .stepPrice(10.0)
                                                .createdAt(now.minusHours(5))
                                                .endAt(now.plusHours(2))
                                                .autoExtend(true)
                                                .bidCount(Long.valueOf(3))
                                                .currentBidder(9L)
                                                .orderCreated(false)
                                                .sentEmail(false)
                                                .build(),

                                // ================= PRODUCT 5 =================
                                Product.builder()
                                                .sellerId(10L)
                                                .name("Limited Edition Sneakers")
                                                .thumbnailUrl("https://example.com/thumb/sneakers.jpg")
                                                .description("Limited edition sneakers, never worn.")
                                                .parentCategoryId(1L)
                                                .parentCategoryName("Fashion")
                                                .categoryId(6L)
                                                .categoryName("Shoes")
                                                .startingPrice(120.0)
                                                .currentPrice(210.0)
                                                .buyNowPrice(300.0)
                                                .stepPrice(10.0)
                                                .createdAt(now.minusDays(4))
                                                .endAt(now.minusHours(2))
                                                .autoExtend(false)
                                                .bidCount(Long.valueOf(6))
                                                .currentBidder(9L)
                                                .orderCreated(false)
                                                .sentEmail(false)
                                                .build());

                productRepository.saveAll(products);
        }
}
