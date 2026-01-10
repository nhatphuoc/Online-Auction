package com.online_auction.bidding_service.dto.response;

import java.time.LocalDateTime;

import com.online_auction.bidding_service.domain.BiddingHistory;

public record BiddingHistorySearchResponse(
        Long id,
        Long productId,
        Long bidderId,
        String bidderName,
        Double amount,
        String requestId,
        BiddingHistory.BidStatus status,
        String reason,
        LocalDateTime createdAt) {
}
