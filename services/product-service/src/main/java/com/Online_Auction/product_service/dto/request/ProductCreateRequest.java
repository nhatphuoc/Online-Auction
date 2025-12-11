package com.Online_Auction.product_service.dto.request;

import java.time.LocalDateTime;
import java.util.List;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter @Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class ProductCreateRequest {

    @NotBlank
    private String name;

    @NotBlank
    private String thumbnailUrl;

    @Size(min = 3)
    private List<String> images;

    @NotBlank
    private String description;

    // ===== CATEGORY INFO SENT FROM UI =====
    @NotNull
    private Long categoryId;

    @NotBlank
    private String categoryName;

    private Long parentCategoryId;
    private String parentCategoryName;

    private Double startingPrice;
    private Double buyNowPrice;
    private Double stepPrice;

    private LocalDateTime endAt;

    private boolean autoExtend;
}

