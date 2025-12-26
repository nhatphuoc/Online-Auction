package com.Online_Auction.product_service.external.order;

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
public class CreateOrderRequest {
    private Long auction_id;
    private Double final_price;
    private Long seller_id;
    private Long winner_id;
}
