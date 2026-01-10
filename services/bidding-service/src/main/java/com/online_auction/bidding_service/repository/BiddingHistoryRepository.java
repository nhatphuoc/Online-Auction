package com.online_auction.bidding_service.repository;

import java.time.LocalDateTime;
import java.util.List;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.JpaSpecificationExecutor;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import com.online_auction.bidding_service.domain.BiddingHistory;

public interface BiddingHistoryRepository
    extends JpaRepository<BiddingHistory, Long>, JpaSpecificationExecutor<BiddingHistory> {

  @Query(value = """
          SELECT b FROM BiddingHistory b
          WHERE (:productId IS NULL OR b.productId = :productId)
            AND (:bidderId IS NULL OR b.bidderId = :bidderId)
            AND (:status IS NULL OR b.status = :status)
            AND (:requestId IS NULL OR b.requestId = :requestId)
            AND (:from IS NULL OR b.createdAt >= :from)
            AND (:to IS NULL OR b.createdAt <= :to)
      """, countQuery = """
          SELECT COUNT(b) FROM BiddingHistory b
          WHERE (:productId IS NULL OR b.productId = :productId)
            AND (:bidderId IS NULL OR b.bidderId = :bidderId)
            AND (:status IS NULL OR b.status = :status)
            AND (:requestId IS NULL OR b.requestId = :requestId)
            AND (:from IS NULL OR b.createdAt >= :from)
            AND (:to IS NULL OR b.createdAt <= :to)
      """)
  Page<BiddingHistory> search(
      @Param("productId") Long productId,
      @Param("bidderId") Long bidderId,
      @Param("status") BiddingHistory.BidStatus status,
      @Param("requestId") String requestId,
      @Param("from") LocalDateTime from,
      @Param("to") LocalDateTime to,
      Pageable pageable);

  @Query("SELECT b FROM BiddingHistory b WHERE b.bidderId = :userId ORDER BY b.createdAt DESC")
  List<BiddingHistory> findAllByBidderId(@Param("userId") Long userId);
}
