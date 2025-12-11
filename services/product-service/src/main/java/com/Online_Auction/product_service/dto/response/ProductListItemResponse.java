package com.Online_Auction.product_service.dto.response;

import java.time.LocalDateTime;

public record ProductListItemResponse(
        Long id,
        String thumbnailUrl,
        String name,
        Double currentPrice,
        Long currentBidder,
        Double buyNowPrice,
        LocalDateTime createdAt,
        LocalDateTime endAt,
        Long bidCount,

        // === NEW FIELDS ===
        Long categoryParentId,
        String categoryParentName,
        Long categoryId,
        String categoryName
) {}

