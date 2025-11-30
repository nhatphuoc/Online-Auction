package com.Online_Auction.auth_service.service;

import java.time.LocalDateTime;
import java.time.temporal.ChronoUnit;

import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import com.Online_Auction.auth_service.domain.OtpToken;
import com.Online_Auction.auth_service.domain.User;
import com.Online_Auction.auth_service.dto.request.RegisterRequest;
import com.Online_Auction.auth_service.dto.request.SignInRequest;
import com.Online_Auction.auth_service.repository.OtpTokenRepository;
import com.Online_Auction.auth_service.repository.UserRepository;

import jakarta.transaction.Transactional;

@Service
public class AuthServiceImpl implements AuthService {

    private final UserRepository userRepository;
    private final RestTemplateNotificationService notificationService;
    private final OtpTokenRepository otpTokenRepository;
    private final PasswordEncoder passwordEncoder;

    public AuthServiceImpl(
        UserRepository userRepository,
        RestTemplateNotificationService notificationService,
        OtpTokenRepository otpTokenRepository,
        PasswordEncoder passwordEncoder
    ) {
        this.userRepository = userRepository;
        this.notificationService = notificationService;
        this.otpTokenRepository = otpTokenRepository;
        this.passwordEncoder = passwordEncoder;
    }

    @Override
    @Transactional
    public void register(RegisterRequest request) {
        // 1. Kiểm tra email tồn tại
        if(userRepository.findByEmail(request.getEmail()).isPresent()){
            throw new RuntimeException("Email already registered");
        }

        // 2. Tạo user
        User user = new User();
        user.setFullName(request.getFullName());
        user.setEmail(request.getEmail());
        user.setPassword(this.passwordEncoder.encode( request.getPassword()));
        user.setBirthDay(request.getBirthDay());
        userRepository.save(user);

        // 3. Tạo OTP
        String otpCode = generateOtp();
        OtpToken otpToken = new OtpToken();
        otpToken.setEmail(request.getEmail());
        otpToken.setOtpCode(otpCode);
        otpToken.setExpiredAt(LocalDateTime.now().plus(10, ChronoUnit.MINUTES));
        otpTokenRepository.save(otpToken);
        
        // 4. Gọi NotificationService
        notificationService.sendEmail(
                user.getEmail(),
                "Xác nhận đăng ký",
                "Mã OTP của bạn là: " + otpCode + "\nHạn dùng: 10 phút"
        );
    }

    private String generateOtp() {
        int otp = (int)(Math.random() * 900000) + 100000;
        return String.valueOf(otp);
    }

    @Transactional
    public boolean verifyOtp(String email, String otpCode) {
        OtpToken otpToken = otpTokenRepository.findByEmail(email)
                .orElseThrow(() -> new RuntimeException("OTP not found"));

        if (otpToken.getExpiredAt().isBefore(LocalDateTime.now())) {
            otpTokenRepository.deleteByEmail(email);
            throw new RuntimeException("OTP expired");
        }

        if (!otpToken.getOtpCode().equals(otpCode)) {
            throw new RuntimeException("Invalid OTP");
        }

        // Đánh dấu user emailVerified = true
        User user = userRepository.findByEmail(email)
                .orElseThrow(() -> new RuntimeException("User not found"));
        user.setEmailVerified(true);
        userRepository.save(user);

        // Xóa OTP sau khi xác thực
        otpTokenRepository.deleteByEmail(email);
        return true;
    }

    @Override
    public User authenticate(SignInRequest request) {
        User user = userRepository.findByEmail(request.getEmail())
                .orElseThrow(() -> new RuntimeException("User not found"));
        if (!passwordEncoder.matches(request.getPassword(), user.getPassword())) {
            throw new RuntimeException("Invalid credentials");
        }
        if (!Boolean.TRUE.equals(user.getEmailVerified())) {
            throw new RuntimeException("Email not verified");
        }
        return user;
    }   
}