package com.Online_Auction.product_service.job;

import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Component
@RequiredArgsConstructor
@Slf4j
public class AuctionFinalizeJob {

    private final AuctionFinalizeHandler auctionFinalizeHandler;

    @Scheduled(cron = "0 */15 * * * *")
    public void run() {
        int processed = auctionFinalizeHandler.finalizeExpiredAuctions();
        log.info("Auction finalize job processed {} auctions", processed);
    }
}
