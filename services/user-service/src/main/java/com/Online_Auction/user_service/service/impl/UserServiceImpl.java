package com.Online_Auction.user_service.service.impl;

import java.util.Optional;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.Online_Auction.user_service.domain.User;
import com.Online_Auction.user_service.domain.User.UserRole;
import com.Online_Auction.user_service.dto.request.RegisterRequest;
import com.Online_Auction.user_service.dto.response.StatusResponse;
import com.Online_Auction.user_service.repository.UserRepository;
import com.Online_Auction.user_service.service.UserService;

@Service
public class UserServiceImpl implements UserService {

    @Autowired
    private UserRepository userRepository;

    @Override
    public User findByEmail(String email) {
        Optional<User> optionalUser = userRepository.findByEmail(email);

        if (optionalUser.isPresent()) return optionalUser.get();
        return null;
    }

    @Override
    public boolean register(RegisterRequest request) {    
        User user = new User();
        user.setFullName(request.getFullName());
        user.setEmail(request.getEmail());
        user.setPassword(request.getPassword());
        user.setBirthDay(request.getBirthDay());
        user.setRole(UserRole.BIDDER);
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
        return new StatusResponse(
            true,
            "Successfully verify email"
        );
    }
    
}
