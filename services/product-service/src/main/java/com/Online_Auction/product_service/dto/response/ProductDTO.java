package com.Online_Auction.product_service.dto.response;

import java.time.LocalDateTime;
import java.util.List;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class ProductDTO {

    private Long id;

    // BASIC INFO
    private String name;
    private String thumbnailUrl;
    private List<String> images;
    private String description;

    // CATEGORY
    private Long categoryId;

    // PRICING
    private Double startingPrice;
    private Double currentPrice;
    private Double buyNowPrice;
    private Double stepPrice;

    // AUCTION
    private LocalDateTime createdAt;
    private LocalDateTime endAt;
    private boolean autoExtend;

    // EXTEND CONFIG
    private Integer extendThresholdMinutes;
    private Integer extendDurationMinutes;

    // SELLER & BIDDER INFO FROM USER-SERVICE
    private Long sellerId;
    private SimpleUserInfo sellerInfo; // CALL USER-SERVICE
    private SimpleUserInfo highestBidder; // FROM BIDDING-SERVICE
}
