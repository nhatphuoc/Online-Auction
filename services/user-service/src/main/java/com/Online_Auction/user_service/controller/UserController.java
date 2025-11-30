package com.Online_Auction.user_service.controller;

import org.springframework.web.bind.annotation.RestController;

import com.Online_Auction.user_service.domain.User;
import com.Online_Auction.user_service.dto.request.RegisterRequest;
import com.Online_Auction.user_service.dto.response.StatusResponse;
import com.Online_Auction.user_service.dto.response.UserResponse;
import com.Online_Auction.user_service.mapper.UserMapper;
import com.Online_Auction.user_service.service.UserService;

import org.springframework.web.bind.annotation.RequestMapping;

import java.util.Map;
import java.util.Objects;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;

@RestController
@RequestMapping("/api/users")
public class UserController {

    @Autowired
    private UserService userService;
    
    @GetMapping
    public UserResponse findByEmail(@RequestParam String email) {
        User user = userService.findByEmail(email);
        if (Objects.isNull(user))
            return null;
        return UserMapper.toUserResponse(userService.findByEmail(email));
    }

    @PostMapping
    public ResponseEntity<StatusResponse> registerUser(@RequestBody RegisterRequest entity) {
        boolean registerSuccess = userService.register(entity);
        if (!registerSuccess) {
            return ResponseEntity.badRequest().body(new StatusResponse(
                false, 
                "Fail to register user")
            );
        }
        return ResponseEntity.ok().body(new StatusResponse(
            registerSuccess, 
            "Successfully register user"
        ));
    }

    @PostMapping("/verify-email")
    public ResponseEntity<StatusResponse> verifyEmail(@RequestBody Map<String, String> payload) {
        String email = payload.get("email");
        return ResponseEntity.ok().body(userService.verifyEmail(email));
    }
    
    @DeleteMapping
    public ResponseEntity<StatusResponse> deleteUserByEmail(@RequestBody Map<String, String> payload) {
        return ResponseEntity.ok().body(userService.deleteUserByEmail(payload.get("email")));
    }
}
