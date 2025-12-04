package com.Online_Auction.auth_service.service;

import com.Online_Auction.auth_service.dto.request.GoogleTokenRequest;
import com.Online_Auction.auth_service.dto.request.RegisterUserRequest;
import com.Online_Auction.auth_service.dto.request.SignInRequest;
import com.Online_Auction.auth_service.external.response.SimpleUserResponse;
import com.Online_Auction.auth_service.external.response.StatusResponse;

public interface AuthService {
    void register(RegisterUserRequest request);
    StatusResponse verifyOtp(String email, String otpCode);
    SimpleUserResponse authenticate(SignInRequest request);
    SimpleUserResponse loginWithGoogle(GoogleTokenRequest request);
}