package com.Online_Auction.auth_service.controller;

import java.io.IOException;
import java.security.GeneralSecurityException;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;
import com.Online_Auction.auth_service.config.jwt.JwtUtils;
import com.Online_Auction.auth_service.dto.request.GoogleTokenRequest;
import com.Online_Auction.auth_service.dto.request.RegisterUserRequest;
import com.Online_Auction.auth_service.dto.request.SignInRequest;
import com.Online_Auction.auth_service.dto.request.ValidateJwtRequest;
import com.Online_Auction.auth_service.dto.request.VerifyOtpRequest;
import com.Online_Auction.auth_service.dto.response.JwtResponse;
import com.Online_Auction.auth_service.dto.response.ValidateJwtResponse;
import com.Online_Auction.auth_service.external.response.ApiResponse;
import com.Online_Auction.auth_service.external.response.SimpleUserResponse;
import com.Online_Auction.auth_service.service.AuthService;

@RestController
@RequestMapping("/auth")
public class AuthController {

    private final AuthService authService;

    public AuthController(AuthService authService, RestTemplate restTemplate) { 
        this.authService = authService;
    }

    @PostMapping("/register")
    public ResponseEntity<ApiResponse<Void>> register(@RequestBody RegisterUserRequest request) {
        ApiResponse<Void> response = authService.register(request);
        return response.isSuccess() ? ResponseEntity.ok().body(response) : ResponseEntity.badRequest().body(response);
    }

    @PostMapping("/verify-otp")
    public ResponseEntity<ApiResponse<Void>> verifyOtp(@RequestBody VerifyOtpRequest request) {
        ApiResponse<Void> response = authService.verifyOtp(request.getEmail(), request.getOtpCode());
        return response.isSuccess() ? ResponseEntity.ok().body(response) : ResponseEntity.badRequest().body(response);
    }

    @PostMapping("/sign-in")
    public ResponseEntity<JwtResponse> signIn(@RequestBody SignInRequest request) {
        SimpleUserResponse user = authService.authenticate(request);

        if (user == null) {
            return ResponseEntity.badRequest().body(new JwtResponse(false,"",""));
        };

        JwtResponse jwtResponse = new JwtResponse(
            true,
            JwtUtils.generateAccessToken(user.getId(), user.getEmail(), user.getUserRole()),
            JwtUtils.generateRefreshToken(user.getId())
        );
        return ResponseEntity.ok(jwtResponse);
    }

    @PostMapping("/validate-jwt")
    public ResponseEntity<ValidateJwtResponse> validateJwt(@RequestBody ValidateJwtRequest request) {
        boolean valid = JwtUtils.validateToken(request.getToken());

        ValidateJwtResponse response = new ValidateJwtResponse(valid);

        if (!valid) {
            return ResponseEntity.status(401).body(response);
        }

        return ResponseEntity.ok(response);
    }

    @PostMapping("/sign-in/google")
    public ResponseEntity<JwtResponse> loginWithGoogle(@RequestBody GoogleTokenRequest request) throws GeneralSecurityException, IOException {
        SimpleUserResponse user = authService.loginWithGoogle(request);

        if (user == null) {
            return ResponseEntity.badRequest().body(
                new JwtResponse(
                    false, 
                    "", 
                    ""
                )
            );
        }
        JwtResponse jwtResponse = new JwtResponse(
            true,
            JwtUtils.generateAccessToken(user.getId(), user.getEmail(), user.getUserRole()),
            JwtUtils.generateRefreshToken(user.getId())
        );
        return ResponseEntity.ok(jwtResponse);
    }

}
