package com.Online_Auction.product_service.service;

import java.time.LocalDateTime;
import java.util.Objects;

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

        // ====== 2. Kiểm tra trạng thái ======
        if (product.getStatus() != Product.ProductStatus.ACTIVE)
            return ApiResponse.fail("Product is not active");

        // ====== 3. Kiểm tra thời gian kết thúc ======
        if (LocalDateTime.now().isAfter(product.getEndAt()))
            return ApiResponse.fail("Auction has already ended");

        Double bidAmount = req.getAmount();

        if (Objects.isNull(product.getCurrentPrice())) {
            if (bidAmount >= product.getStartingPrice()) {
                product.setCurrentBidder(req.getBidderId());
                product.setCurrentPrice(bidAmount);
                productRepository.save(product);

                ProductBidSuccessData data = new ProductBidSuccessData(
                    bidAmount,
                    null
                );

                return ApiResponse.success(data, "Bid placed successfully");
            }
            return ApiResponse.fail("Bid amount must be equal or higher than starting price");
        }

        // ====== 4. Giá phải > giá hiện tại ======
        if (bidAmount <= product.getCurrentPrice())
            return ApiResponse.fail("Bid amount must be higher than current price");

        // ====== 5. Giá phải theo đúng stepPrice ======
        if (product.getStepPrice() != null) {
            double diff = bidAmount - product.getCurrentPrice();
            if (diff % product.getStepPrice() != 0) {
                return ApiResponse.fail("Bid must follow step price: " + product.getStepPrice());
            }
        }

        // ====== 6. Lưu bidder cũ ======
        Long previousHighestBidder = null;
        // giả sử có bảng lưu highest bidder → bỏ qua
        // ở đây có thể return null
        // hoặc bạn có thể thêm trường highestBidder vào Product

        // ====== 7. Cập nhật giá ======
        product.setCurrentPrice(bidAmount);
        productRepository.save(product);

        // ====== 8. Trả về data success ======
        ProductBidSuccessData data = new ProductBidSuccessData(
                bidAmount,
                previousHighestBidder
        );

        return ApiResponse.success(data, "Bid placed successfully");
    }
}
