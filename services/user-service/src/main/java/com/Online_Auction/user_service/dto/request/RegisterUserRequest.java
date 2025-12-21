package com.Online_Auction.user_service.dto.request;

import java.time.LocalDate;

public record RegisterUserRequest(
    String fullName,
    String email,
    String password,
    LocalDate birthDay,
    Boolean emailVerified
) {}