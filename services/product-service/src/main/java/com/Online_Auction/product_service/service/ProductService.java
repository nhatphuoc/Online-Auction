package com.Online_Auction.product_service.service;

import lombok.RequiredArgsConstructor;

import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.server.ResponseStatusException;

import com.Online_Auction.product_service.client.RestTemplateUserServiceClient;
import com.Online_Auction.product_service.domain.Product;
import com.Online_Auction.product_service.domain.Product.ProductStatus;
import com.Online_Auction.product_service.dto.request.ProductCreateRequest;
import com.Online_Auction.product_service.dto.request.ProductUpdateRequest;
import com.Online_Auction.product_service.dto.response.ProductDTO;
import com.Online_Auction.product_service.dto.response.SimpleUserInfo;
import com.Online_Auction.product_service.external.SimpleUserResponse;
import com.Online_Auction.product_service.mapper.ProductMapper;
import com.Online_Auction.product_service.repository.ProductRepository;

import java.time.LocalDateTime;
import java.util.List;

@Service
@RequiredArgsConstructor
public class ProductService {

    private final ProductRepository productRepository;
    private final ProductMapper productMapper;
    private final RestTemplateUserServiceClient restTemplateUserServiceClient;

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
                .startingPrice(request.getStartingPrice())
                .buyNowPrice(request.getBuyNowPrice())
                .stepPrice(request.getStepPrice())
                .createdAt(LocalDateTime.now())
                .endAt(request.getEndAt())
                .autoExtend(request.isAutoExtend())
                .sellerId(sellerId)
                .status(ProductStatus.ACTIVE)
                .build();

        productRepository.save(product);

        SimpleUserInfo sellerInfo = this.getSimpleUserInfoById(sellerId);
        SimpleUserInfo highestBidder = null;

        return productMapper.toProductDTO(product, sellerInfo, highestBidder);
    }

    // =================================
    // GET PRODUCT DETAIL (ALL USER)
    // =================================
    @Transactional(readOnly = true)
    public ProductDTO getProductDetail(Long productId) {
        Product product = productRepository.findById(productId)
                .orElseThrow(() -> new IllegalArgumentException("Product not found"));

        SimpleUserInfo sellerInfo = this.getSimpleUserInfoById(product.getSellerId());
        SimpleUserInfo highestBidder = null;        // TODO: call bidding-service

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
        SimpleUserInfo highestBidder = null;        // TODO

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

    private SimpleUserInfo getSimpleUserInfoById(long id) {
        SimpleUserResponse userResponse = restTemplateUserServiceClient.getUserById(id);
        if (userResponse == null) {
            throw new ResponseStatusException(
                HttpStatus.NOT_FOUND,
                "User not found with id = " + id
            );
        }
        SimpleUserInfo userInfo = new SimpleUserInfo();
        userInfo.setId(userInfo.getId());
        userInfo.setEmail(userInfo.getEmail());
        userInfo.setFullName(userInfo.getFullName());
        userInfo.setUserRole(userInfo.getUserRole());
        return userInfo;
    }
}
