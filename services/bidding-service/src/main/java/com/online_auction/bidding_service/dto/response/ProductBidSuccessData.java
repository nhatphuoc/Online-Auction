package com.online_auction.bidding_service.dto.response;

import lombok.*;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class ProductBidSuccessData {
    private Double newHighest;
    private Long previousHighestBidder; // nullable
}