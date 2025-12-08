package com.online_auction.bidding_service.controller;

import lombok.RequiredArgsConstructor;

import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

import com.online_auction.bidding_service.config.security.UserPrincipal;
import com.online_auction.bidding_service.dto.request.BidRequest;
import com.online_auction.bidding_service.dto.response.ApiResponse;
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
            @RequestHeader("X-User-Token") String userJwt
    ) {
        Authentication auth = SecurityContextHolder.getContext().getAuthentication();
        UserPrincipal user = (UserPrincipal) auth.getPrincipal();

        ApiResponse<?> response = bidService.placeBid(
                req.getProductId(),
                user.getUserId(),
                req.getAmount(),
                req.getRequestId(),
                userJwt
        );

        return ResponseEntity
                .status(response.isSuccess() ? 200 : 400)
                .body(response);
    }
}