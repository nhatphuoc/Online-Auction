package com.Online_Auction.user_service.service;

import com.Online_Auction.user_service.domain.User;
import com.Online_Auction.user_service.dto.request.RegisterUserRequest;
import com.Online_Auction.user_service.dto.request.SignInRequest;
import com.Online_Auction.user_service.dto.response.SimpleUserResponse;
import com.Online_Auction.user_service.dto.response.StatusResponse;

public interface UserService {
    User findByEmail(String email);
    boolean register(RegisterUserRequest registerRequest);
    StatusResponse verifyEmail(String email);
    StatusResponse deleteUserByEmail(String email);
    User findById(long id);
    SimpleUserResponse authenticateUser(SignInRequest request);
    User getCurrentUser();
}