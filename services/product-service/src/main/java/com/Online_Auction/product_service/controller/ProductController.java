package com.Online_Auction.product_service.controller;

import lombok.RequiredArgsConstructor;

import org.springframework.security.core.Authentication;
import org.springframework.data.domain.Page;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

import com.Online_Auction.product_service.config.security.UserPrincipal;
import com.Online_Auction.product_service.dto.request.ProductCreateRequest;
import com.Online_Auction.product_service.dto.request.ProductUpdateRequest;
import com.Online_Auction.product_service.dto.request.RenameParentCategoryRequest;
import com.Online_Auction.product_service.dto.request.UpdateCategoryRequest;
import com.Online_Auction.product_service.dto.response.ApiResponse;
import com.Online_Auction.product_service.dto.response.BatchUpdateResult;
import com.Online_Auction.product_service.dto.response.ProductDTO;
import com.Online_Auction.product_service.dto.response.ProductListItemResponse;
import com.Online_Auction.product_service.external.ProductBidRequest;
import com.Online_Auction.product_service.service.ProductBidService;
import com.Online_Auction.product_service.service.ProductService;

import jakarta.validation.Valid;

import java.util.List;

@RestController
@RequestMapping("/api/products")
@RequiredArgsConstructor
public class ProductController {

    private final ProductService productService;
    private final ProductBidService productBidService;

    // =================================
    // SELLER: CREATE PRODUCT
    // =================================
    @PreAuthorize("hasRole('SELLER')")
    @PostMapping
    public ResponseEntity<ProductDTO> createProduct(
            @Valid @RequestBody ProductCreateRequest request) {
        // Lấy userId từ SecurityContext
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        UserPrincipal principal = (UserPrincipal) authentication.getPrincipal();
        Long sellerId = principal.getUserId();

        ProductDTO product = productService.createProduct(sellerId, request);
        return ResponseEntity.ok(product);
    }

    // =================================
    // GET PRODUCT DETAIL (ALL USERS)
    // =================================
    @GetMapping("/{productId}")
    public ResponseEntity<ProductDTO> getProductDetail(@PathVariable Long productId) {
        ProductDTO product = productService.getProductDetail(productId);
        return ResponseEntity.ok(product);
    }

    // =================================
    // SELLER: UPDATE DESCRIPTION (APPEND)
    // =================================
    @PreAuthorize("hasRole('SELLER')")
    @PatchMapping("/{productId}/description")
    public ResponseEntity<ProductDTO> updateDescription(
            @PathVariable Long productId,
            @Valid @RequestBody ProductUpdateRequest request) {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        UserPrincipal principal = (UserPrincipal) authentication.getPrincipal();
        Long sellerId = principal.getUserId();
        ProductDTO updated = productService.updateProductDescription(sellerId, productId, request);
        return ResponseEntity.ok(updated);
    }

    // =================================
    // LIST PRODUCT BY SELLER
    // =================================
    @GetMapping("/seller/{sellerId}")
    public ResponseEntity<List<ProductDTO>> listProductsBySeller(@PathVariable Long sellerId) {
        List<ProductDTO> products = productService.listProductsBySeller(sellerId);
        return ResponseEntity.ok(products);
    }

    // =================================
    // INTERNAL USAGE: BID
    // =================================
    @PostMapping("/{id}/bids")
    public ResponseEntity<?> placeBid(
            @PathVariable Long id,
            @RequestBody ProductBidRequest request) {
        ApiResponse<?> response = productBidService.placeBid(id, request);

        if (!response.isSuccess())
            return ResponseEntity.badRequest().body(response);

        return ResponseEntity.ok(response);
    }

    // =================================
    // HOMEPAGE
    // =================================
    @GetMapping("/top-ending")
    public ApiResponse<List<ProductListItemResponse>> topEnding() {
        return ApiResponse.success(
                productService.topEndingSoon(),
                "Successfully fetching top5 ending-soon products");
    }

    @GetMapping("/top-most-bids")
    public ApiResponse<List<ProductListItemResponse>> topMostBids() {
        return ApiResponse.success(
                productService.topMostBids(),
                "Successfully fetching top5 most-bids products");
    }

    @GetMapping("/top-highest-price")
    public ApiResponse<List<ProductListItemResponse>> topHighestPrice() {
        return ApiResponse.success(
                productService.topHighestPrice(),
                "Successfully fetching top5 highest-price products");
    }

    // =================================
    // SEARCH + FILTER
    // =================================
    @GetMapping("/search")
    public ApiResponse<Page<ProductListItemResponse>> searchProducts(
            @RequestParam(required = false) String query,
            @RequestParam(required = false) Long parentCategoryId,
            @RequestParam(required = false) Long categoryId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int pageSize) {
        Page<ProductListItemResponse> result = productService.searchProducts(
                query,
                parentCategoryId,
                categoryId,
                page,
                pageSize);

        return ApiResponse.success(result, "Query success");
    }

    /**
     * ================= UPDATE CATEGORY =================
     *
     * Use-case:
     * - Category-service gọi endpoint này khi:
     * + ĐỔI TÊN category (categoryName)
     * + DI CHUYỂN category sang parent khác
     * + Đồng bộ lại snapshot category info trong Product
     *
     * Những gì endpoint này làm:
     * - Update tất cả Product có categoryId = {categoryId}
     * - Cập nhật các field:
     * + categoryName
     * + parentCategoryId
     * + parentCategoryName
     *
     * Những gì endpoint này KHÔNG làm:
     * - KHÔNG validate category tồn tại hay không (category-service chịu trách
     * nhiệm)
     * - KHÔNG update product thuộc category con khác
     *
     * Ví dụ use-case:
     * - "Laptop Gaming" → đổi tên thành "Gaming Laptop"
     * - Category 5 được move từ "Electronics" sang "Tech Devices"
     *
     * Endpoint:
     * PUT /api/internal/products/categories/{categoryId}
     */

    @PutMapping("/categories/{categoryId}")
    public ResponseEntity<ApiResponse<BatchUpdateResult>> updateCategory(
            @PathVariable Long categoryId,
            @RequestBody UpdateCategoryRequest request) {
        int updated = productService.updateCategory(categoryId, request);

        return ResponseEntity.ok(
                ApiResponse.success(
                        new BatchUpdateResult(updated),
                        "Category updated successfully"));
    }

    /**
     * ================= RENAME PARENT CATEGORY =================
     *
     * Use-case:
     * - Category-service gọi endpoint này khi:
     * + CHỈ ĐỔI TÊN parent category
     *
     * Những gì endpoint này làm:
     * - Update tất cả Product có parentCategoryId = {parentCategoryId}
     * - CHỈ cập nhật:
     * + parentCategoryName
     *
     * Những gì endpoint này KHÔNG làm:
     * - KHÔNG thay đổi categoryId
     * - KHÔNG thay đổi categoryName
     * - KHÔNG move category
     *
     * Ví dụ use-case:
     * - Parent category "Electronics" → đổi tên thành "Tech Devices"
     * - Tất cả product thuộc các category con đều được đồng bộ tên parent
     *
     * Endpoint:
     * PUT /api/internal/products/parent-categories/{parentCategoryId}/rename
     */

    @PutMapping("/parent-categories/{parentCategoryId}/rename")
    public ResponseEntity<ApiResponse<BatchUpdateResult>> renameParentCategory(
            @PathVariable Long parentCategoryId,
            @RequestBody RenameParentCategoryRequest request) {
        int updated = productService.renameParentCategory(
                parentCategoryId,
                request.getParentCategoryName());

        return ResponseEntity.ok(
                ApiResponse.success(
                        new BatchUpdateResult(updated),
                        "Parent category renamed successfully"));
    }
}
