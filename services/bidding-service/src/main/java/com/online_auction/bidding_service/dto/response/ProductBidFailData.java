package com.online_auction.bidding_service.dto.response;

import lombok.*;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class ProductBidFailData {
    private Double currentHighest;
    private String reason;
}