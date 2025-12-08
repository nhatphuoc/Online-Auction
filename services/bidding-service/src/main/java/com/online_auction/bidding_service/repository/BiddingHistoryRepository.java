package com.online_auction.bidding_service.repository;

import org.springframework.data.jpa.repository.JpaRepository;

import com.online_auction.bidding_service.domain.BiddingHistory;

public interface BiddingHistoryRepository extends JpaRepository<BiddingHistory, Long> {
}
