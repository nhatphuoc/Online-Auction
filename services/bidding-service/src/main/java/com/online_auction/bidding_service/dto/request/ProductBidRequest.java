package com.online_auction.bidding_service.dto.request;

import lombok.*;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class ProductBidRequest {
    private Long bidderId;
    private Double amount;
    private String requestId;
}
