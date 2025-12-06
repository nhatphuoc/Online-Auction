package com.Online_Auction.auth_service.external.response;

import java.time.LocalDate;

import com.Online_Auction.auth_service.external.response.SimpleUserResponse.UserRole;

public record UserProfileResponse(
    Long id,
    String fullName,
    String email,
    LocalDate birthDay,
    UserRole role,
    int totalNumberReviews,
    int totalNumberGoodReviews,
    Boolean emailVerified,
    Boolean isSellerRequestSent
) {}
