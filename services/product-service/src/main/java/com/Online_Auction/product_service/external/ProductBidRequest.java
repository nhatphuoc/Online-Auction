package com.Online_Auction.product_service.external;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class ProductBidRequest {
    private Long bidderId;
    private Double amount;
    private String requestId;
}
