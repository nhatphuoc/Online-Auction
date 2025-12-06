package com.Online_Auction.user_service.service.impl;

import java.util.Optional;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import com.Online_Auction.user_service.config.security.UserPrincipal;
import com.Online_Auction.user_service.domain.User;
import com.Online_Auction.user_service.domain.User.UserRole;
import com.Online_Auction.user_service.dto.request.RegisterUserRequest;
import com.Online_Auction.user_service.dto.request.SignInRequest;
import com.Online_Auction.user_service.dto.response.SimpleUserResponse;
import com.Online_Auction.user_service.dto.response.StatusResponse;
import com.Online_Auction.user_service.mapper.UserMapper;
import com.Online_Auction.user_service.repository.UserRepository;
import com.Online_Auction.user_service.service.UserService;

import io.micrometer.common.util.StringUtils;

@Service
public class UserServiceImpl implements UserService {

    @Autowired
    private UserRepository userRepository;

    private PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();

    @Override
    public User findByEmail(String email) {
        Optional<User> optionalUser = userRepository.findByEmail(email);

        if (optionalUser.isPresent()) return optionalUser.get();
        return null;
    }

    @Override
    public boolean register(RegisterUserRequest request) {
        Optional<User> existing = userRepository.findByEmail(request.email());
        if (existing.isPresent()) return false;
        
        User user = new User();
        user.setFullName(request.fullName());
        user.setEmail(request.email());

        if (StringUtils.isNotBlank(request.prodiver())) {
            String hashPassword = passwordEncoder.encode(request.password());
            user.setPassword(hashPassword);
        } else {
            user.setPassword(null);
        }
        
        user.setBirthDay(request.birthDay());
        user.setEmailVerified(request.emailVerified());
        user.setRole(UserRole.ROLE_BIDDER);
        userRepository.save(user);
        return true;
    }

    @Override
    public StatusResponse verifyEmail(String email) {
        Optional<User> optional = userRepository.findByEmail(email);
        if (!optional.isPresent()) {
            return new StatusResponse(
                false, 
                "Fail to verify email, email not exists"
            );
        }
        User user = optional.get();
        user.setEmailVerified(true);
        userRepository.save(user);
        return new StatusResponse(
            true,
            "Successfully verify email"
        );
    }

    @Override
    public StatusResponse deleteUserByEmail(String email) {
        userRepository.deleteByEmail(email);
        return new StatusResponse(
            true,
            "Successfully delete user by email"
        );
    }

    @Override
    public User findById(long id) {
        Optional<User> optional = userRepository.findById(id);
        if (!optional.isPresent())
            return null;
        return optional.get();
    }

    @Override
    public SimpleUserResponse authenticateUser(SignInRequest request) {
        Optional<User> user = userRepository.findByEmail(request.getEmail());
        if (!user.isPresent()) {
            return null;
        }
        return UserMapper.toSimpleUserResponse(user.get());
    }

    public User getCurrentUser() {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();

        if (authentication == null || !authentication.isAuthenticated()) {
            return null;
        }

        UserPrincipal principal = (UserPrincipal) authentication.getPrincipal();
        return userRepository.findById(principal.getUserId()).orElse(null);
    }

    
}
