package com.Online_Auction.product_service.controller;

import lombok.RequiredArgsConstructor;

import org.springframework.security.core.Authentication;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

import com.Online_Auction.product_service.config.security.UserPrincipal;
import com.Online_Auction.product_service.dto.FavoriteDTO;
import com.Online_Auction.product_service.dto.ProductDTO;
import com.Online_Auction.product_service.dto.QuestionDTO;
import com.Online_Auction.product_service.dto.request.AnswerCreateRequest;
import com.Online_Auction.product_service.dto.request.ProductCreateRequest;
import com.Online_Auction.product_service.dto.request.ProductUpdateRequest;
import com.Online_Auction.product_service.dto.request.QuestionCreateRequest;
import com.Online_Auction.product_service.service.FavoriteService;
import com.Online_Auction.product_service.service.ProductService;
import com.Online_Auction.product_service.service.QuestionService;

import jakarta.validation.Valid;

import java.util.List;

@RestController
@RequestMapping("/api/products")
@RequiredArgsConstructor
public class ProductController {

    private final ProductService productService;
    private final QuestionService questionService;
    private final FavoriteService favoriteService;

    // =================================
    // SELLER: CREATE PRODUCT
    // =================================
    @PreAuthorize("hasRole('SELLER')")
    @PostMapping
    public ResponseEntity<ProductDTO> createProduct(
            @Valid @RequestBody ProductCreateRequest request
    ) {
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
            @Valid @RequestBody ProductUpdateRequest request
    ) {
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
    // BIDDER: ASK QUESTION
    // =================================
    @PreAuthorize("hasRole('BIDDER')")
    @PostMapping("/{productId}/questions")
    public ResponseEntity<QuestionDTO> askQuestion(
            @PathVariable Long productId,
            @Valid @RequestBody QuestionCreateRequest request
    ) {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        UserPrincipal principal = (UserPrincipal) authentication.getPrincipal();
        Long userId = principal.getUserId();
        QuestionDTO question = questionService.askQuestion(userId, productId, request);
        return ResponseEntity.ok(question);
    }

    // =================================
    // SELLER: ANSWER QUESTION
    // =================================
    @PreAuthorize("hasRole('SELLER')")
    @PostMapping("/questions/{questionId}/answer")
    public ResponseEntity<QuestionDTO> answerQuestion(
            @PathVariable Long questionId,
            @Valid @RequestBody AnswerCreateRequest request
    ) {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        UserPrincipal principal = (UserPrincipal) authentication.getPrincipal();
        Long sellerId = principal.getUserId();
        QuestionDTO question = questionService.answerQuestion(sellerId, questionId, request);
        return ResponseEntity.ok(question);
    }

    // =================================
    // FAVORITE: ADD
    // =================================
    @PreAuthorize("isAuthenticated()")
    @PostMapping("/{productId}/favorite")
    public ResponseEntity<Void> addFavorite(@PathVariable Long productId) {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        UserPrincipal principal = (UserPrincipal) authentication.getPrincipal();

        Long userId = principal.getUserId(); // Lấy userId từ principal
        favoriteService.addFavorite(userId, productId);
        return ResponseEntity.ok().build();
    }

    // =================================
    // FAVORITE: REMOVE
    // =================================
    @PreAuthorize("isAuthenticated()")
    @DeleteMapping("/{productId}/favorite")
    public ResponseEntity<Void> removeFavorite(
            @PathVariable Long productId
    ) {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        UserPrincipal principal = (UserPrincipal) authentication.getPrincipal();

        Long userId = principal.getUserId(); // Lấy userId từ principal
        favoriteService.removeFavorite(userId, productId);
        return ResponseEntity.ok().build();
    }

    // =================================
    // LIST FAVORITE PRODUCTS OF USER
    // =================================
    @PreAuthorize("isAuthenticated()")
    @GetMapping("/favorites")
    public ResponseEntity<List<FavoriteDTO>> listFavorites() {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        UserPrincipal principal = (UserPrincipal) authentication.getPrincipal();

        Long userId = principal.getUserId(); // Lấy userId từ principal
        List<FavoriteDTO> favorites = favoriteService.listFavorites(userId);
        return ResponseEntity.ok(favorites);
    }
}
