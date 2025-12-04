package com.Online_Auction.auth_service.dto.request;

import java.time.LocalDate;

public record RegisterUserRequest(
    String fullName,
    String email,
    String password,
    LocalDate birthDay,
    Boolean emailVerified
) {}