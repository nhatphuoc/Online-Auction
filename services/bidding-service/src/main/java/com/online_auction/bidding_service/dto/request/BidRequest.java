package com.online_auction.bidding_service.dto.request;

import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotNull;
import lombok.*;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class BidRequest {
    @NotNull
    private Long productId;

    @NotNull
    @Min(0)
    private Double amount;

    // Optional client side request id for idempotency/tracing
    private String requestId;
}