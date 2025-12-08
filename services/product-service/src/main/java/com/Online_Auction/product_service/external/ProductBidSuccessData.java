package com.Online_Auction.product_service.external;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class ProductBidSuccessData {
    private Double newHighest;
    private Long previousHighestBidder;
}
