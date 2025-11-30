package com.Online_Auction.auth_service.service;

import com.Online_Auction.auth_service.dto.request.RegisterRequest;
import com.Online_Auction.auth_service.dto.request.SignInRequest;
import com.Online_Auction.auth_service.external.response.StatusResponse;
import com.Online_Auction.auth_service.external.response.UserProfileResponse;

public interface AuthService {
    void register(RegisterRequest request);
    StatusResponse verifyOtp(String email, String otpCode);
    UserProfileResponse authenticate(SignInRequest request);
}