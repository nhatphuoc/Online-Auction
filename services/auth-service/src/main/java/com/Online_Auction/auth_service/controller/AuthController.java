package com.Online_Auction.auth_service.controller;

import java.util.Map;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.Online_Auction.auth_service.config.jwt.JwtUtils;
import com.Online_Auction.auth_service.domain.User;
import com.Online_Auction.auth_service.dto.request.RegisterRequest;
import com.Online_Auction.auth_service.dto.request.SignInRequest;
import com.Online_Auction.auth_service.dto.request.ValidateJwtRequest;
import com.Online_Auction.auth_service.dto.request.VerifyOtpRequest;
import com.Online_Auction.auth_service.dto.response.JwtResponse;
import com.Online_Auction.auth_service.dto.response.ValidateJwtResponse;
import com.Online_Auction.auth_service.service.AuthService;

@RestController
@RequestMapping("/auth")
public class AuthController {
    private final AuthService authService;

    public AuthController(AuthService authService) { this.authService = authService; }

    @PostMapping("/register")
    public ResponseEntity<Map<String,Object>> register(@RequestBody RegisterRequest request) {
        authService.register(request);
        return ResponseEntity.ok(Map.of(
            "success", true, 
            "message", "Successfully register"
        ));
    }

    @PostMapping("/verify-otp")
    public ResponseEntity<Map<String, Object>> verifyOtp(@RequestBody VerifyOtpRequest request) {
        boolean success = authService.verifyOtp(request.getEmail(), request.getOtpCode());
        return ResponseEntity.ok(Map.of(
            "success", success,
            "message", "OTP verified successfully"
        ));
    }

    @PostMapping("/sign-in")
    public ResponseEntity<JwtResponse> signIn(@RequestBody SignInRequest request) {
        User user = authService.authenticate(request);

        JwtResponse jwtResponse = new JwtResponse(
                JwtUtils.generateAccessToken(user),
                JwtUtils.generateRefreshToken(user)
        );
        return ResponseEntity.ok(jwtResponse);
    }

    @PostMapping("/validate-jwt")
    public ResponseEntity<ValidateJwtResponse> validateJwt(@RequestBody ValidateJwtRequest request) {
        boolean valid = JwtUtils.validateToken(request.getToken());

        ValidateJwtResponse response = new ValidateJwtResponse(valid);

        if (!valid) {
            // trả 401 nếu token không hợp lệ
            return ResponseEntity.status(401).body(response);
        }

        return ResponseEntity.ok(response);
    }
}
