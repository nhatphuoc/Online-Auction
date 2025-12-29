package com.online_auction.bidding_service.specs;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;

import org.springframework.data.jpa.domain.Specification;

import com.online_auction.bidding_service.domain.BiddingHistory;
import com.online_auction.bidding_service.domain.BiddingHistory.BidStatus;

import jakarta.persistence.criteria.Predicate;

public class BiddingHistorySpecs {
    public static Specification<BiddingHistory> search(
            Long productId,
            Long bidderId,
            BidStatus status,
            String requestId,
            LocalDateTime from,
            LocalDateTime to) {

        return (root, query, cb) -> {
            List<Predicate> p = new ArrayList<>();

            if (productId != null)
                p.add(cb.equal(root.get("productId"), productId));
            if (bidderId != null)
                p.add(cb.equal(root.get("bidderId"), bidderId));
            if (status != null)
                p.add(cb.equal(root.get("status"), status));
            if (requestId != null)
                p.add(cb.equal(root.get("requestId"), requestId));
            if (from != null)
                p.add(cb.greaterThanOrEqualTo(root.get("createdAt"), from));
            if (to != null)
                p.add(cb.lessThanOrEqualTo(root.get("createdAt"), to));

            return cb.and(p.toArray(new Predicate[0]));
        };
    }

}
