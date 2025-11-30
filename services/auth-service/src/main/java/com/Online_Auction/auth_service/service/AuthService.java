package com.Online_Auction.auth_service.service;

import com.Online_Auction.auth_service.domain.User;
import com.Online_Auction.auth_service.dto.request.RegisterRequest;
import com.Online_Auction.auth_service.dto.request.SignInRequest;

public interface AuthService {
    void register(RegisterRequest request);
    boolean verifyOtp(String email, String otpCode);
    User authenticate(SignInRequest request);
}