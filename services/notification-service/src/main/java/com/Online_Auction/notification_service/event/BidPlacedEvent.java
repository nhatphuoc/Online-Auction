package com.Online_Auction.notification_service.event;

import java.io.Serializable;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
@JsonIgnoreProperties(ignoreUnknown = true)
public class BidPlacedEvent implements Serializable {
    private Long productId;
    private Long bidderId;
    private Double amount;
    private Long previousHighestBidder;
    private String requestId;
}
