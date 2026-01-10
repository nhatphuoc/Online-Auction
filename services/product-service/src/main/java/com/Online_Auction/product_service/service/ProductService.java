package com.Online_Auction.product_service.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.data.jpa.domain.Specification;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.server.ResponseStatusException;

import com.Online_Auction.product_service.client.RestTemplateNotificationServiceClient;
import com.Online_Auction.product_service.client.RestTemplateOrderServiceClient;
import com.Online_Auction.product_service.client.RestTemplateUserServiceClient;
import com.Online_Auction.product_service.controller.ProductController.ProductSort;
import com.Online_Auction.product_service.domain.Product;
import com.Online_Auction.product_service.dto.request.ProductCreateRequest;
import com.Online_Auction.product_service.dto.request.ProductUpdateRequest;
import com.Online_Auction.product_service.dto.request.UpdateCategoryRequest;
import com.Online_Auction.product_service.dto.response.BuyNowResponse;
import com.Online_Auction.product_service.dto.response.ProductDTO;
import com.Online_Auction.product_service.dto.response.ProductListItemResponse;
import com.Online_Auction.product_service.dto.response.SimpleUserInfo;
import com.Online_Auction.product_service.external.SimpleUserResponse;
import com.Online_Auction.product_service.external.notification.EmailNotificationRequest;
import com.Online_Auction.product_service.external.order.CreateOrderRequest;
import com.Online_Auction.product_service.mapper.ProductMapper;
import com.Online_Auction.product_service.repository.ProductRepository;
import com.Online_Auction.product_service.specs.ProductSpecs;
import com.fasterxml.jackson.databind.ObjectMapper;

import java.time.LocalDateTime;
import java.util.List;

@Service
@RequiredArgsConstructor
@Slf4j
public class ProductService {

    private ObjectMapper om = new ObjectMapper();

    private final ProductRepository productRepository;
    private final ProductMapper productMapper;
    private final RestTemplateUserServiceClient restTemplateUserServiceClient;
    private final RestTemplateOrderServiceClient orderClient;
    private final RestTemplateNotificationServiceClient notificationServiceClient;

    // =================================
    // CREATE PRODUCT (SELLER)
    // =================================
    @Transactional
    public ProductDTO createProduct(Long sellerId, ProductCreateRequest request) {

        Product product = Product.builder()
                .name(request.getName())
                .thumbnailUrl(request.getThumbnailUrl())
                .images(request.getImages())
                .description(request.getDescription())

                .categoryId(request.getCategoryId())
                .categoryName(request.getCategoryName())
                .parentCategoryId(request.getParentCategoryId())
                .parentCategoryName(request.getParentCategoryName())

                .startingPrice(request.getStartingPrice())
                .buyNowPrice(request.getBuyNowPrice())
                .stepPrice(request.getStepPrice())
                .createdAt(LocalDateTime.now())
                .endAt(request.getEndAt())
                .autoExtend(request.isAutoExtend())
                .sellerId(sellerId)
                .build();

        productRepository.save(product);

        SimpleUserInfo sellerInfo = this.getSimpleUserInfoById(sellerId);

        return productMapper.toProductDTO(product, sellerInfo, null);
    }

    // =================================
    // GET PRODUCT DETAIL (ALL USER)
    // =================================
    @Transactional(readOnly = true)
    public ProductDTO getProductDetail(Long productId) {
        Product product = productRepository.findById(productId)
                .orElseThrow(() -> new IllegalArgumentException("Product not found"));

        SimpleUserInfo sellerInfo = this.getSimpleUserInfoById(product.getSellerId());
        SimpleUserInfo highestBidder = this.getSimpleUserInfoById(product.getCurrentBidder());

        return productMapper.toProductDTO(product, sellerInfo, highestBidder);
    }

    // =================================
    // UPDATE PRODUCT DESCRIPTION (SELLER)
    // =================================
    @Transactional
    public ProductDTO updateProductDescription(Long sellerId, Long productId, ProductUpdateRequest request) {
        Product product = productRepository.findById(productId)
                .orElseThrow(() -> new IllegalArgumentException("Product not found"));

        if (!product.getSellerId().equals(sellerId)) {
            throw new IllegalArgumentException("You are not the seller of this product");
        }

        // Bổ sung mô tả (append)
        String newDescription = product.getDescription() + "\n" + request.getAdditionalDescription();
        product.setDescription(newDescription);

        productRepository.save(product);

        SimpleUserInfo sellerInfo = new SimpleUserInfo(); // TODO
        SimpleUserInfo highestBidder = null; // TODO

        return productMapper.toProductDTO(product, sellerInfo, highestBidder);
    }

    // =================================
    // LIST PRODUCT BY SELLER
    // =================================
    @Transactional(readOnly = true)
    public List<ProductDTO> listProductsBySeller(Long sellerId) {
        List<Product> products = productRepository.findBySellerId(sellerId);

        return products.stream()
                .map(p -> productMapper.toProductDTO(p, new SimpleUserInfo(), null))
                .toList();
    }

    // =================================
    // DELETE PRODUCT (Optional)
    // =================================
    @Transactional
    public void deleteProduct(Long sellerId, Long productId) {
        Product product = productRepository.findById(productId)
                .orElseThrow(() -> new IllegalArgumentException("Product not found"));

        if (!product.getSellerId().equals(sellerId)) {
            throw new IllegalArgumentException("You are not the seller of this product");
        }

        productRepository.delete(product);
    }

    // =================================
    // HOMEPAGE
    // =================================
    public List<ProductListItemResponse> topEndingSoon() {
        return productRepository.findTop5EndingSoon(
                LocalDateTime.now(),
                PageRequest.of(0, 5))
                .stream()
                .map(ProductMapper::toListItem)
                .toList();
    }

    public List<ProductListItemResponse> topMostBids() {
        return productRepository.findTop5MostBids(
                LocalDateTime.now(),
                PageRequest.of(0, 5))
                .stream()
                .map(ProductMapper::toListItem)
                .toList();
    }

    public List<ProductListItemResponse> topHighestPrice() {
        return productRepository.findTop5HighestPrice(
                LocalDateTime.now(),
                PageRequest.of(0, 5))
                .stream()
                .map(ProductMapper::toListItem)
                .toList();
    }

    // =================================
    // SEARCH + FILTER
    // =================================
    public Page<ProductListItemResponse> searchProducts(
            String query,
            Long categoryId,
            int page,
            int pageSize,
            ProductSort sort) {

        Sort sortSpec = switch (sort) {
            case PRICE_ASC -> Sort.by("currentPrice").ascending();
            case PRICE_DESC -> Sort.by("currentPrice").descending();
            case BID_COUNT_DESC -> Sort.by("bidCount").descending();
            case NEWEST -> Sort.by("createdAt").descending();
        };

        Pageable pageable = PageRequest.of(page, pageSize, sortSpec);

        Specification<Product> spec = Specification
                .where(ProductSpecs.hasCategory(categoryId))
                .and(ProductSpecs.hasNamePrefix(query));

        return productRepository.findAll(spec, pageable)
                .map(p -> new ProductListItemResponse(
                        p.getId(),
                        p.getThumbnailUrl(),
                        p.getName(),
                        p.getCurrentPrice(),
                        p.getBuyNowPrice(),
                        p.getCreatedAt(),
                        p.getEndAt(),
                        p.getBidCount(),
                        p.getParentCategoryId(),
                        p.getParentCategoryName(),
                        p.getCategoryId(),
                        p.getCategoryName()));
    }

    public Page<ProductListItemResponse> getWonProducts(Long bidderId, int page, int pageSize) {
        LocalDateTime now = LocalDateTime.now();

        Pageable pageable = PageRequest.of(page, pageSize, Sort.by("endAt").descending());

        Page<Product> products = productRepository.findByHighestBidderAndEnded(bidderId, now, pageable);

        return products.map(p -> new ProductListItemResponse(
                p.getId(),
                p.getThumbnailUrl(),
                p.getName(),
                p.getCurrentPrice(),
                p.getBuyNowPrice(),
                p.getCreatedAt(),
                p.getEndAt(),
                p.getBidCount(),
                p.getParentCategoryId(),
                p.getParentCategoryName(),
                p.getCategoryId(),
                p.getCategoryName()));
    }

    // =================================
    // UTILITY FUNCTIONS
    // =================================
    private SimpleUserInfo getSimpleUserInfoById(long id) {
        SimpleUserResponse userResponse = restTemplateUserServiceClient.getUserById(id);
        try {
            log.info("Raw: {}", om.writeValueAsString(userResponse));
        } catch (Exception e) {
            // TODO: handle exception
        }

        if (userResponse == null) {
            throw new ResponseStatusException(
                    HttpStatus.NOT_FOUND,
                    "User not found with id = " + id);
        }
        SimpleUserInfo userInfo = new SimpleUserInfo();
        userInfo.setId(userResponse.getId());
        userInfo.setEmail(userResponse.getEmail());
        userInfo.setFullName(userResponse.getFullName());
        userInfo.setUserRole(userResponse.getUserRole());
        return userInfo;
    }

    /* ================= UPDATE CATEGORY ================= */

    @Transactional
    public int updateCategory(Long categoryId, UpdateCategoryRequest request) {

        if (categoryId == null) {
            throw new IllegalArgumentException("categoryId must not be null");
        }

        return productRepository.updateByCategoryId(
                categoryId,
                request.getCategoryName(),
                request.getParentCategoryId(),
                request.getParentCategoryName());
    }

    /* ================= RENAME PARENT CATEGORY ================= */

    @Transactional
    public int renameParentCategory(Long parentCategoryId, String parentCategoryName) {

        if (parentCategoryId == null) {
            throw new IllegalArgumentException("parentCategoryId must not be null");
        }

        if (parentCategoryName == null || parentCategoryName.isBlank()) {
            throw new IllegalArgumentException("parentCategoryName must not be blank");
        }

        return productRepository.updateParentCategoryName(
                parentCategoryId,
                parentCategoryName);
    }

    /* ================= BUY NOW ================= */
    @Transactional
    public BuyNowResponse buyNow(Long productId, Long buyerId) {

        Product product = productRepository.findByIdForUpdate(productId)
                .orElseThrow(() -> new RuntimeException("Product not found"));

        LocalDateTime now = LocalDateTime.now();

        // ===== VALIDATION =====
        if (product.isOrderCreated()) {
            throw new RuntimeException("Product already sold");
        }

        if (product.getBuyNowPrice() == null) {
            throw new RuntimeException("Buy now not available");
        }

        if (product.getEndAt().isBefore(now)) {
            throw new RuntimeException("Auction already ended");
        }

        // ===== MARK AS WON =====
        markAuctionAsWon(product, buyerId, product.getBuyNowPrice());

        productRepository.save(product);

        // ===== SIDE EFFECTS =====
        handleWinningAuction(product);

        return new BuyNowResponse(
                product.getId(),
                product.getCurrentPrice(),
                buyerId,
                product.getEndAt());
    }

    private void markAuctionAsWon(Product product, Long winnerId, Double finalPrice) {
        LocalDateTime now = LocalDateTime.now();

        product.setCurrentBidder(winnerId);
        product.setCurrentPrice(finalPrice);
        product.setEndAt(now);
    }

    private void handleWinningAuction(Product product) {

        if (!product.isOrderCreated()) {
            orderClient.createOrder(
                    CreateOrderRequest.builder()
                            .auction_id(product.getId())
                            .final_price(product.getCurrentPrice())
                            .seller_id(product.getSellerId())
                            .winner_id(product.getCurrentBidder())
                            .build());

            product.setOrderCreated(true);
            productRepository.save(product);
        }

        if (!product.isSentEmail()) {
            var seller = restTemplateUserServiceClient.getUserById(product.getSellerId());
            var winner = restTemplateUserServiceClient.getUserById(product.getCurrentBidder());

            notificationServiceClient.sendEmail(
                    EmailNotificationRequest.builder()
                            .to(seller.getEmail())
                            .subject("Auction ended with a winner")
                            .body("Your product \"" + product.getName() + "\" has been sold.")
                            .build());

            notificationServiceClient.sendEmail(
                    EmailNotificationRequest.builder()
                            .to(winner.getEmail())
                            .subject("You won the auction!")
                            .body("You won \"" + product.getName() + "\".")
                            .build());

            product.setSentEmail(true);
            productRepository.save(product);
        }
    }
}
