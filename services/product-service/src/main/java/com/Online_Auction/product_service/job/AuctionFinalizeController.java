package com.Online_Auction.product_service.job;

import java.util.Map;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import lombok.RequiredArgsConstructor;

@RestController
@RequiredArgsConstructor
@RequestMapping("/api/products/internal/auctions")
public class AuctionFinalizeController {

    private final AuctionFinalizeHandler auctionFinalizeHandler;

    @PostMapping("/finalize")
    public ResponseEntity<?> finalizeManually() {
        int processed = auctionFinalizeHandler.finalizeExpiredAuctions();
        return ResponseEntity.ok(
                Map.of(
                        "processed", processed,
                        "message", "Auction finalization triggered successfully"));
    }
}
