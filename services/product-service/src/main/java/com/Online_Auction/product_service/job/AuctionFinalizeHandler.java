package com.Online_Auction.product_service.job;

import java.util.List;

import org.springframework.stereotype.Service;

import com.Online_Auction.product_service.client.RestTemplateNotificationServiceClient;
import com.Online_Auction.product_service.client.RestTemplateOrderServiceClient;
import com.Online_Auction.product_service.client.RestTemplateUserServiceClient;
import com.Online_Auction.product_service.domain.Product;
import com.Online_Auction.product_service.external.notification.EmailNotificationRequest;
import com.Online_Auction.product_service.external.order.CreateOrderRequest;
import com.Online_Auction.product_service.repository.ProductRepository;

import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Service
@RequiredArgsConstructor
@Slf4j
public class AuctionFinalizeHandler {

    private final ProductRepository productRepository;
    private final RestTemplateUserServiceClient userClient;
    private final RestTemplateOrderServiceClient orderClient;
    private final RestTemplateNotificationServiceClient notificationClient;

    @Transactional
    public int finalizeExpiredAuctions() {
        List<Product> products = productRepository.findExpiredAuctionsForProcessing();

        if (products.isEmpty()) {
            return 0;
        }

        int processedProducts = 0;

        for (Product product : products) {
            try {
                if (product.getBidCount() > 0 && product.getCurrentBidder() != null) {
                    handleWinningAuction(product);
                } else {
                    notifyNoBid(product);
                }

                product.setOrderCreated(true);
                product.setSentEmail(true);
                productRepository.save(product);
                processedProducts++;
            } catch (Exception ex) {
                log.error("Failed to process auction {}", product.getId(), ex);
            }
        }

        return processedProducts;
    }

    private void handleWinningAuction(Product product) {
        var seller = userClient.getUserById(product.getSellerId());
        var winner = userClient.getUserById(product.getCurrentBidder());

        if (!product.isOrderCreated()) {
            orderClient.createOrder(
                    CreateOrderRequest.builder()
                            .auction_id(product.getId())
                            .final_price(product.getCurrentPrice())
                            .seller_id(product.getSellerId())
                            .winner_id(product.getCurrentBidder())
                            .build());
        }

        if (!product.isSentEmail()) {
            notificationClient.sendEmail(
                    EmailNotificationRequest.builder()
                            .to(seller.getEmail())
                            .subject("Auction ended with a winner")
                            .body("Your product \"" + product.getName() + "\" has been sold.")
                            .build());

            notificationClient.sendEmail(
                    EmailNotificationRequest.builder()
                            .to(winner.getEmail())
                            .subject("You won the auction!")
                            .body("You won \"" + product.getName() + "\".")
                            .build());
        }
    }

    private void notifyNoBid(Product product) {
        var seller = userClient.getUserById(product.getSellerId());

        if (!product.isSentEmail()) {
            notificationClient.sendEmail(
                    EmailNotificationRequest.builder()
                            .to(seller.getEmail())
                            .subject("Auction ended without bids")
                            .body("Your product \"" + product.getName() + "\" ended with no bids.")
                            .build());
        }
    }
}
