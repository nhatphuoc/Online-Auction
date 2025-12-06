package com.Online_Auction.user_service.controller;


import org.springframework.web.bind.annotation.RestController;

import com.Online_Auction.user_service.domain.User;
import com.Online_Auction.user_service.dto.request.RegisterUserRequest;
import com.Online_Auction.user_service.dto.request.SignInRequest;
import com.Online_Auction.user_service.dto.response.ApiResponse;
import com.Online_Auction.user_service.dto.response.SimpleUserResponse;
import com.Online_Auction.user_service.dto.response.StatusResponse;
import com.Online_Auction.user_service.dto.response.UserProfileResponse;
import com.Online_Auction.user_service.mapper.UserMapper;
import com.Online_Auction.user_service.service.UserService;

import lombok.RequiredArgsConstructor;

import org.springframework.web.bind.annotation.RequestMapping;

import java.util.Map;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;

@RestController
@RequestMapping("/api/users")
@RequiredArgsConstructor
public class UserController {

    private final UserService userService;

    @GetMapping("/simple")
    public ResponseEntity<ApiResponse<SimpleUserResponse>> getSimpleUserByEmail(@RequestParam String email) {
        User user = userService.findByEmail(email);
        if (user == null)
            return ResponseEntity.badRequest().body(ApiResponse.fail("User not found"));

        return ResponseEntity.ok(
                ApiResponse.success(UserMapper.toSimpleUserResponse(user), "User fetched successfully")
        );
    }

    @PostMapping
    public ResponseEntity<ApiResponse<Void>> registerUser(@RequestBody RegisterUserRequest entity) {
        boolean registerSuccess = userService.register(entity);

        if (!registerSuccess) {
            return ResponseEntity.badRequest()
                    .body(ApiResponse.fail("Fail to register user, email is already registered"));
        }

        System.out.println("Success");
        return ResponseEntity.ok(ApiResponse.success(null, "Successfully registered user"));
    }

    @PostMapping("/verify-email")
    public ResponseEntity<ApiResponse<Void>> verifyEmail(@RequestBody Map<String, String> payload) {
        String email = payload.get("email");
        StatusResponse status = userService.verifyEmail(email);

        if (!status.isSuccess())
            return ResponseEntity.badRequest().body(ApiResponse.fail(status.getMessage()));

        return ResponseEntity.ok(ApiResponse.success(null, status.getMessage()));
    }

    @DeleteMapping
    public ResponseEntity<ApiResponse<Void>> deleteUserByEmail(@RequestBody Map<String, String> payload) {
        StatusResponse status = userService.deleteUserByEmail(payload.get("email"));

        if (!status.isSuccess())
            return ResponseEntity.badRequest().body(ApiResponse.fail(status.getMessage()));

        return ResponseEntity.ok(ApiResponse.success(null, status.getMessage()));
    }

    @PostMapping("/authenticate")
    public ResponseEntity<ApiResponse<SimpleUserResponse>> authenticateUser(@RequestBody SignInRequest request) {
        SimpleUserResponse response = userService.authenticateUser(request);

        return ResponseEntity.ok(ApiResponse.success(
                response,
                "Authentication successful"
        ));
    }

    /**
     * GET /api/users/profile/me
     */
    @GetMapping("/profile/me")
    public ResponseEntity<ApiResponse<UserProfileResponse>> getUserProfile() {
        User user = userService.getCurrentUser();

        if (user == null)
            return ResponseEntity.badRequest().body(ApiResponse.fail("User not authenticated"));

        return ResponseEntity.ok(
                ApiResponse.success(UserMapper.toUserProfileResponse(user), "Profile retrieved")
        );
    }
}
