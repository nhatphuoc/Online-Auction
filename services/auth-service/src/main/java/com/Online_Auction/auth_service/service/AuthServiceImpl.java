package com.Online_Auction.auth_service.service;

import java.time.LocalDateTime;
import java.time.temporal.ChronoUnit;

import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

import com.Online_Auction.auth_service.domain.OtpToken;
import com.Online_Auction.auth_service.dto.request.GoogleTokenRequest;
import com.Online_Auction.auth_service.dto.request.RegisterUserRequest;
import com.Online_Auction.auth_service.dto.request.SignInRequest;
import com.Online_Auction.auth_service.external.client.UserServiceClient;
import com.Online_Auction.auth_service.external.response.ApiResponse;
import com.Online_Auction.auth_service.external.response.SimpleUserResponse;
import com.Online_Auction.auth_service.repository.OtpTokenRepository;
import com.google.api.client.googleapis.auth.oauth2.GoogleIdToken;
import com.google.api.client.googleapis.auth.oauth2.GoogleIdTokenVerifier;

import jakarta.transaction.Transactional;

@Service
public class AuthServiceImpl implements AuthService {

    private final RestTemplateNotificationService notificationService;
    private final OtpTokenRepository otpTokenRepository;
    private final UserServiceClient userServiceClient;
    private final GoogleIdTokenVerifier googleIdTokenVerifier;

    public AuthServiceImpl(
            RestTemplateNotificationService notificationService,
            OtpTokenRepository otpTokenRepository,
            UserServiceClient userServiceClient,
            GoogleIdTokenVerifier googleIdTokenVerifier) {
        this.notificationService = notificationService;
        this.otpTokenRepository = otpTokenRepository;
        this.userServiceClient = userServiceClient;
        this.googleIdTokenVerifier = googleIdTokenVerifier;
    }

    @Override
    @Transactional
    public ApiResponse<Void> register(RegisterUserRequest request) {
        // 1. Tạo user
        ApiResponse<Void> responseRegisterUser = userServiceClient.registerUser(request);
        if (!responseRegisterUser.isSuccess()) {
            return new ApiResponse<Void>(
                    false,
                    responseRegisterUser.getMessage(),
                    null);
        }

        // 2. Tạo OTP
        String otpCode = generateOtp();
        OtpToken otpToken = otpTokenRepository.findByEmail(request.email())
                .orElse(new OtpToken());
        otpToken.setEmail(request.email());
        otpToken.setOtpCode(otpCode);
        otpToken.setExpiredAt(LocalDateTime.now().plus(10, ChronoUnit.MINUTES));
        otpTokenRepository.save(otpToken);

        // 4. Gọi NotificationService
        if (notificationService != null) {
            notificationService.sendEmail(
                    request.email(),
                    "Xác nhận đăng ký",
                    "Mã OTP của bạn là: " + otpCode + "\nHạn dùng: 10 phút");
        }

        return responseRegisterUser;
    }

    private String generateOtp() {
        int otp = (int) (Math.random() * 900000) + 100000;
        return String.valueOf(otp);
    }

    @Transactional
    public ApiResponse<Void> verifyOtp(String email, String otpCode) {
        OtpToken otpToken = otpTokenRepository.findByEmail(email)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.BAD_REQUEST, "OTP not found"));

        if (otpToken.getExpiredAt().isBefore(LocalDateTime.now())) {
            otpTokenRepository.deleteByEmail(email);
            userServiceClient.deleteUserByEmail(email);
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "OTP expired");
        }

        if (!otpToken.getOtpCode().equals(otpCode)) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Invalid OTP");
        }

        // Đánh dấu user emailVerified = true
        ApiResponse<Void> statusResponse = userServiceClient.verifyEmail(email);
        if (!statusResponse.isSuccess()) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, statusResponse.getMessage());
        }

        // Xóa OTP sau khi xác thực
        otpTokenRepository.deleteByEmail(email);
        return statusResponse;
    }

    @Override
    public SimpleUserResponse authenticate(SignInRequest request) {
        ApiResponse<SimpleUserResponse> authenticateResponse = userServiceClient.authenticateUser(request);
        if (!authenticateResponse.isSuccess()) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, authenticateResponse.getMessage());
        }
        return authenticateResponse.getData();
    }

    @Override
    public SimpleUserResponse loginWithGoogle(GoogleTokenRequest request) {

        GoogleIdToken idToken;
        try {
            idToken = googleIdTokenVerifier.verify(request.idToken());
        } catch (Exception e) {
            throw new ResponseStatusException(HttpStatus.UNAUTHORIZED, "Invalid Google ID Token");
        }

        if (idToken == null) {
            throw new ResponseStatusException(HttpStatus.UNAUTHORIZED, "Invalid Google ID Token");
        }

        GoogleIdToken.Payload payload = idToken.getPayload();

        String email = payload.getEmail();
        String name = (String) payload.get("name");

        // 1. Get user from user-service
        ApiResponse<SimpleUserResponse> getSimpleUserResponse = userServiceClient.getUserByEmail(email);

        if (!getSimpleUserResponse.isSuccess()) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, getSimpleUserResponse.getMessage());
        }

        // 2. Automatically register user
        RegisterUserRequest req = new RegisterUserRequest(
                name,
                email,
                null,
                null,
                true);

        ApiResponse<Void> status = userServiceClient.registerUser(req);
        if (!status.isSuccess()) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Failed to auto-register user");
        }

        ApiResponse<SimpleUserResponse> savedUser = userServiceClient.getUserByEmail(email);
        return savedUser.getData();
    }
}