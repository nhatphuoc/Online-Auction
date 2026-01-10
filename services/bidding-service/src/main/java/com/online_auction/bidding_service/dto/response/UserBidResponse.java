package com.online_auction.bidding_service.dto.response;

import java.time.LocalDateTime;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class UserBidResponse {
    private Long bidId;
    private Double bidAmount;
    private String bidStatus;
    private String reason;
    private LocalDateTime bidCreatedAt;

    // Product info
    private Long productId;
    private String productName;
    private String thumbnailUrl;
    private Double currentPrice;
    private Double buyNowPrice;
    private LocalDateTime endAt;
    private boolean autoExtend;
    private Long currentBidder;
}
