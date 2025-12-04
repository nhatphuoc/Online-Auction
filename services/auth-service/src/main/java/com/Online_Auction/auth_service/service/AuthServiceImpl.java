package com.Online_Auction.auth_service.service;

import java.time.LocalDateTime;
import java.time.temporal.ChronoUnit;
import java.util.Objects;

import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

import com.Online_Auction.auth_service.domain.OtpToken;
import com.Online_Auction.auth_service.dto.request.GoogleTokenRequest;
import com.Online_Auction.auth_service.dto.request.RegisterUserRequest;
import com.Online_Auction.auth_service.dto.request.SignInRequest;
import com.Online_Auction.auth_service.external.response.SimpleUserResponse;
import com.Online_Auction.auth_service.external.response.StatusResponse;
import com.Online_Auction.auth_service.repository.OtpTokenRepository;
import com.google.api.client.googleapis.auth.oauth2.GoogleIdToken;
import com.google.api.client.googleapis.auth.oauth2.GoogleIdTokenVerifier;

import jakarta.transaction.Transactional;

@Service
public class AuthServiceImpl implements AuthService {

    private final RestTemplateNotificationService notificationService;
    private final OtpTokenRepository otpTokenRepository;
    private final RestTemplateUserService restTemplateUserService;
    private final GoogleIdTokenVerifier googleIdTokenVerifier;

    public AuthServiceImpl(
        RestTemplateNotificationService notificationService,
        OtpTokenRepository otpTokenRepository,
        RestTemplateUserService restTemplateUserService,
        GoogleIdTokenVerifier googleIdTokenVerifier
    ) {
        this.notificationService = notificationService;
        this.otpTokenRepository = otpTokenRepository;
        this.restTemplateUserService = restTemplateUserService;
        this.googleIdTokenVerifier = googleIdTokenVerifier;
    }

    @Override
    @Transactional
    public void register(RegisterUserRequest request) {
        // 1. Kiểm tra email tồn tại
        SimpleUserResponse userResponse = restTemplateUserService.getUserByEmail(request.email());
        if (!Objects.isNull(userResponse)) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Email already registered");
        }

        // 2. Tạo user
        StatusResponse response = restTemplateUserService.registerUser(request);
        if (!response.isSuccess()) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Fail to register new user, unknown reason");
        }

        // 3. Tạo OTP
        String otpCode = generateOtp();
        OtpToken otpToken = new OtpToken();
        otpToken.setEmail(request.email());
        otpToken.setOtpCode(otpCode);
        otpToken.setExpiredAt(LocalDateTime.now().plus(10, ChronoUnit.MINUTES));
        otpTokenRepository.save(otpToken);
        
        // 4. Gọi NotificationService
        if (notificationService != null) {
                notificationService.sendEmail(
                request.email(),
                "Xác nhận đăng ký",
                "Mã OTP của bạn là: " + otpCode + "\nHạn dùng: 10 phút"
            );
        }
    }

    private String generateOtp() {
        int otp = (int)(Math.random() * 900000) + 100000;
        return String.valueOf(otp);
    }

    @Transactional
    public StatusResponse verifyOtp(String email, String otpCode) {
        OtpToken otpToken = otpTokenRepository.findByEmail(email)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.BAD_REQUEST, "OTP not found"));

        if (otpToken.getExpiredAt().isBefore(LocalDateTime.now())) {
            otpTokenRepository.deleteByEmail(email);
            restTemplateUserService.deleteUserByEmail(email);
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "OTP expired");
        }

        if (!otpToken.getOtpCode().equals(otpCode)) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Invalid OTP");
        }

        // Đánh dấu user emailVerified = true
        StatusResponse statusResponse = restTemplateUserService.verifyEmail(email);
        if (!statusResponse.isSuccess()) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, statusResponse.getMessage());
        }

        // Xóa OTP sau khi xác thực
        otpTokenRepository.deleteByEmail(email);
        return statusResponse;
    }

    @Override
    public SimpleUserResponse authenticate(SignInRequest request) {
        return restTemplateUserService.authenticateUser(request);
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
        SimpleUserResponse user = restTemplateUserService.getUserByEmail(email);


        // 2. If no user → register automatically
        if (user == null) {
            RegisterUserRequest req = new RegisterUserRequest(
                    name,
                    email,
                    null,
                    null,
                    true
            );

            StatusResponse status = restTemplateUserService.registerUser(req);
            if (!status.isSuccess()) {
                throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Failed to auto-register user");
            }

            user = restTemplateUserService.getUserByEmail(email);
            return user;
        }
        return null;
    }
}