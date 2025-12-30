package com.Online_Auction.product_service.dto.response;

import java.time.LocalDateTime;

import lombok.AllArgsConstructor;
import lombok.Data;

@Data
@AllArgsConstructor
public class BuyNowResponse {
    private Long productId;
    private Double finalPrice;
    private Long buyerId;
    private LocalDateTime endAt;
}
