package com.Online_Auction.product_service.domain;

import java.time.LocalDateTime;
import java.util.List;

import jakarta.persistence.CollectionTable;
import jakarta.persistence.Column;
import jakarta.persistence.ElementCollection;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.Table;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Entity
@Table(name = "products")
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class Product {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    // ===== SELLER INFO =====
    private Long sellerId; // từ JWT

    // ===== BASIC INFORMATION =====
    private String name;

    @Column(columnDefinition = "TEXT")
    private String thumbnailUrl; // ảnh đại diện

    @ElementCollection
    @CollectionTable(name = "product_images", joinColumns = @JoinColumn(name = "product_id"))
    @Column(name = "image_url")
    private List<String> images; // ảnh phụ (ít nhất 3 ảnh)

    @Column(columnDefinition = "TEXT")
    private String description; // mô tả (WYSIWYG)

    // ===== CATEGORY =====
    private Long parentCategoryId;
    private String parentCategoryName;

    private Long categoryId;
    private String categoryName;


    // ===== PRICING =====
    private Double startingPrice;

    @Builder.Default
    private Double currentPrice = null;
    private Double buyNowPrice;
    private Double stepPrice;     // bước giá

    // ===== AUCTION TIME =====
    private LocalDateTime createdAt;
    private LocalDateTime endAt;

    private boolean autoExtend; // có hỗ trợ tự động gia hạn không

    // ===== OTHER =====
    @Enumerated(EnumType.STRING)
    private ProductStatus status;

    @Builder.Default
    private Long bidCount = 0L;

    private Long currentBidder;

    public enum ProductStatus {
        ACTIVE,      // đang đấu giá
        FINISHED,    // đã kết thúc
        PENDING,     // đang chờ duyệt
        REJECTED
    }
}
