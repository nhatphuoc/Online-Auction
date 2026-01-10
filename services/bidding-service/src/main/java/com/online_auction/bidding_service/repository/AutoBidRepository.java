package com.online_auction.bidding_service.repository;

import java.util.List;
import java.util.Optional;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import com.online_auction.bidding_service.domain.AutoBid;

@Repository
public interface AutoBidRepository extends JpaRepository<AutoBid, Long> {

    Optional<AutoBid> findByProductIdAndBidderId(Long productId, Long bidderId);

    List<AutoBid> findByProductIdAndActiveTrueOrderByMaxAmountDesc(Long productId);
}
