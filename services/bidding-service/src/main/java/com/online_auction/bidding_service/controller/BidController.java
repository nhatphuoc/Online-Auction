package com.online_auction.bidding_service.controller;

import lombok.RequiredArgsConstructor;

import java.time.LocalDateTime;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.format.annotation.DateTimeFormat;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

import com.online_auction.bidding_service.config.security.UserPrincipal;
import com.online_auction.bidding_service.domain.BiddingHistory;
import com.online_auction.bidding_service.dto.request.BidRequest;
import com.online_auction.bidding_service.dto.response.ApiResponse;
import com.online_auction.bidding_service.dto.response.BiddingHistorySearchResponse;
import com.online_auction.bidding_service.service.BidService;

@RestController
@RequestMapping("/bids")
@RequiredArgsConstructor
public class BidController {

        private final BidService bidService;

        @PostMapping
        @PreAuthorize("hasAnyRole('BIDDER', 'SELLER')")
        public ResponseEntity<?> placeBid(
                        @RequestBody BidRequest req,
                        @RequestHeader("X-User-Token") String userJwt) {
                Authentication auth = SecurityContextHolder.getContext().getAuthentication();
                UserPrincipal user = (UserPrincipal) auth.getPrincipal();

                ApiResponse<?> response = bidService.placeBid(
                                req.getProductId(),
                                user.getUserId(),
                                req.getAmount(),
                                req.getRequestId(),
                                userJwt);

                return ResponseEntity
                                .status(response.isSuccess() ? 200 : 400)
                                .body(response);
        }

        @GetMapping("/search")
        @PreAuthorize("hasAnyRole('ADMIN', 'SELLER', 'BIDDER')")
        public Page<BiddingHistorySearchResponse> search(
                        @RequestParam(required = false) Long productId,
                        @RequestParam(required = false) Long bidderId,
                        @RequestParam(required = false) BiddingHistory.BidStatus status,
                        @RequestParam(required = false) String requestId,
                        @RequestParam(required = false) @DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME) LocalDateTime from,
                        @RequestParam(required = false) @DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME) LocalDateTime to,
                        @RequestParam(defaultValue = "0") int page,
                        @RequestParam(defaultValue = "10") int size) {
                Pageable pageable = PageRequest.of(
                                page,
                                size,
                                Sort.by("createdAt").descending());

                return this.bidService.search(
                                productId,
                                bidderId,
                                status,
                                requestId,
                                from,
                                to,
                                pageable);
        }
}