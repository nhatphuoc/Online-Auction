package com.online_auction.bidding_service.dto.request;

import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
public class AutoBidRegisterRequest {

    @NotNull(message = "productId is required")
    private Long productId;

    @NotNull(message = "maxAmount is required")
    @Positive(message = "maxAmount must be greater than 0")
    private Double maxAmount;
}
