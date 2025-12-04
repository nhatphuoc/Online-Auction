package com.Online_Auction.auth_service.controller;

import java.io.IOException;
import java.security.GeneralSecurityException;
import java.util.Map;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.Online_Auction.auth_service.config.jwt.JwtUtils;
import com.Online_Auction.auth_service.dto.request.GoogleTokenRequest;
import com.Online_Auction.auth_service.dto.request.RegisterUserRequest;
import com.Online_Auction.auth_service.dto.request.SignInRequest;
import com.Online_Auction.auth_service.dto.request.ValidateJwtRequest;
import com.Online_Auction.auth_service.dto.request.VerifyOtpRequest;
import com.Online_Auction.auth_service.dto.response.JwtResponse;
import com.Online_Auction.auth_service.dto.response.ValidateJwtResponse;
import com.Online_Auction.auth_service.external.response.SimpleUserResponse;
import com.Online_Auction.auth_service.external.response.StatusResponse;
import com.Online_Auction.auth_service.service.AuthService;

@RestController
@RequestMapping("/auth")
public class AuthController {
    private final AuthService authService;

    public AuthController(AuthService authService) { this.authService = authService; }

    @PostMapping("/register")
    public ResponseEntity<Map<String,Object>> register(@RequestBody RegisterUserRequest request) {
        System.out.println("RegisterRequest: " + request);
        authService.register(request);
        return ResponseEntity.ok(Map.of(
            "success", true, 
            "message", "Successfully register"
        ));
    }

    @PostMapping("/verify-otp")
    public ResponseEntity<StatusResponse> verifyOtp(@RequestBody VerifyOtpRequest request) {
        StatusResponse response = authService.verifyOtp(request.getEmail(), request.getOtpCode());
        return ResponseEntity.ok().body(response);
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
        System.out.println("Sign in with google");
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
