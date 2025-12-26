package com.Online_Auction.product_service.service;

import java.time.LocalDateTime;

import org.springframework.stereotype.Service;

import com.Online_Auction.product_service.domain.Product;
import com.Online_Auction.product_service.dto.response.ApiResponse;
import com.Online_Auction.product_service.external.ProductBidRequest;
import com.Online_Auction.product_service.external.ProductBidSuccessData;
import com.Online_Auction.product_service.repository.ProductRepository;

import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;

@Service
@RequiredArgsConstructor
public class ProductBidService {

    private final ProductRepository productRepository;

    @Transactional
    public ApiResponse<?> placeBid(Long productId, ProductBidRequest req) {

        // ====== 1. Lấy product với PESSIMISTIC_LOCK ======
        Product product = productRepository.findByIdForUpdate(productId)
                .orElse(null);

        if (product == null)
            return ApiResponse.fail("Product not found");

        // ====== 2. Kiểm tra thời gian kết thúc ======
        if (LocalDateTime.now().isAfter(product.getEndAt()))
            return ApiResponse.fail("Auction has already ended");

        Double bidAmount = req.getAmount();

        // ====== CASE 1: Chưa có ai đặt giá ======
        if (product.getCurrentPrice() == null) {

            if (bidAmount >= product.getStartingPrice()) {

                product.setCurrentBidder(req.getBidderId());
                product.setCurrentPrice(bidAmount);

                // Tăng số lượt ra giá
                product.setBidCount(product.getBidCount() + 1);

                productRepository.save(product);

                return ApiResponse.success(
                        new ProductBidSuccessData(bidAmount, null),
                        "Bid placed successfully");
            }

            return ApiResponse.fail(
                    "Bid amount must be equal or higher than starting price");
        }

        // ====== CASE 2: Đã có người đặt giá ======
        // 4. Giá phải > giá hiện tại
        if (bidAmount <= product.getCurrentPrice())
            return ApiResponse.fail("Bid amount must be higher than current price");

        // 6. Lưu bidder cũ
        Long previousHighestBidder = product.getCurrentBidder();

        // 7. Cập nhật giá + bidder
        product.setCurrentPrice(bidAmount);
        product.setCurrentBidder(req.getBidderId());

        // Tăng lượt ra giá
        product.setBidCount(product.getBidCount() + 1);

        // 8. Lưu DB
        productRepository.save(product);

        // ====== 8. Trả về data success ======
        ProductBidSuccessData data = new ProductBidSuccessData(
                bidAmount,
                previousHighestBidder);

        return ApiResponse.success(data, "Bid placed successfully");
    }
}
