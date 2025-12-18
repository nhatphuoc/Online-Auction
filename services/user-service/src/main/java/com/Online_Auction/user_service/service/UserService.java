package com.Online_Auction.user_service.service;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;

import com.Online_Auction.user_service.domain.User;
import com.Online_Auction.user_service.dto.request.RegisterUserRequest;
import com.Online_Auction.user_service.dto.request.SignInRequest;
import com.Online_Auction.user_service.dto.response.SimpleUserResponse;
import com.Online_Auction.user_service.dto.response.StatusResponse;
import com.Online_Auction.user_service.dto.response.UserSearchResponse;

public interface UserService {
    User findByEmail(String email);

    boolean register(RegisterUserRequest registerRequest);

    StatusResponse verifyEmail(String email);

    StatusResponse deleteUserByEmail(String email);

    User findById(long id);

    SimpleUserResponse authenticateUser(SignInRequest request);

    User getCurrentUser();

    Page<UserSearchResponse> searchUsers(String keyword, User.UserRole role, Pageable pageable);
}