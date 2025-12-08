package com.online_auction.bidding_service.event;

import java.io.Serializable;

import lombok.*;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class BidPlacedEvent implements Serializable {
    private Long productId;
    private Long bidderId;
    private Double amount;
    private Long previousHighestBidder;
    private String requestId;
}