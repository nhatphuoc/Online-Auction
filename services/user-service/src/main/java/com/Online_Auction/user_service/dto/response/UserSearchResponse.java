package com.Online_Auction.user_service.dto.response;

import java.time.LocalDate;

import com.Online_Auction.user_service.domain.User;

public record UserSearchResponse(
        Long id,
        String fullName,
        String email,
        LocalDate birthDay,
        User.UserRole role,
        int totalNumberReviews,
        int totalNumberGoodReviews,
        Boolean emailVerified,
        Boolean isSellerRequestSent) {
}
