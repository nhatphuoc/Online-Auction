package com.Online_Auction.product_service.external.order;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@JsonIgnoreProperties(ignoreUnknown = true)
public class OrderResponse {
    private Long id;
    private Long auction_id;
    private Long seller_id;
    private Long winner_id;
    private Double final_price;
    private String status;
}
