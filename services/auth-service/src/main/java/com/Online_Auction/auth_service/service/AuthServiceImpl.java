package com.Online_Auction.auth_service.service;

import java.time.LocalDateTime;
import java.time.temporal.ChronoUnit;
import java.util.Objects;

import org.springframework.http.HttpStatus;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

import com.Online_Auction.auth_service.domain.OtpToken;
import com.Online_Auction.auth_service.dto.request.RegisterRequest;
import com.Online_Auction.auth_service.dto.request.SignInRequest;
import com.Online_Auction.auth_service.external.response.StatusResponse;
import com.Online_Auction.auth_service.external.response.UserResponse;
import com.Online_Auction.auth_service.repository.OtpTokenRepository;

import jakarta.transaction.Transactional;

@Service
public class AuthServiceImpl implements AuthService {

    private final RestTemplateNotificationService notificationService;
    private final OtpTokenRepository otpTokenRepository;
    private final PasswordEncoder passwordEncoder;
    private final RestTemplateUserService restTemplateUserService;

    public AuthServiceImpl(
        RestTemplateNotificationService notificationService,
        OtpTokenRepository otpTokenRepository,
        PasswordEncoder passwordEncoder,
        RestTemplateUserService restTemplateUserService
    ) {
        this.notificationService = notificationService;
        this.otpTokenRepository = otpTokenRepository;
        this.passwordEncoder = passwordEncoder;
        this.restTemplateUserService = restTemplateUserService;
    }

    @Override
    @Transactional
    public void register(RegisterRequest request) {
        // 1. Kiểm tra email tồn tại
        UserResponse userResponse = restTemplateUserService.getUserByEmail(request.getEmail());
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
        otpToken.setEmail(request.getEmail());
        otpToken.setOtpCode(otpCode);
        otpToken.setExpiredAt(LocalDateTime.now().plus(10, ChronoUnit.MINUTES));
        otpTokenRepository.save(otpToken);
        
        // 4. Gọi NotificationService
        notificationService.sendEmail(
            request.getEmail(),
            "Xác nhận đăng ký",
            "Mã OTP của bạn là: " + otpCode + "\nHạn dùng: 10 phút"
        );
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
    public UserResponse authenticate(SignInRequest request) {
        UserResponse userResponse = restTemplateUserService.getUserByEmail(request.getEmail());
        if (Objects.isNull(userResponse)) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "User not found");
        }

        if (!passwordEncoder.matches(request.getPassword(), userResponse.getPassword())) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Invalid credentials");
        }
        if (!Boolean.TRUE.equals(userResponse.getEmailVerified())) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Email not verified");
        }
        return userResponse;
    }   
}