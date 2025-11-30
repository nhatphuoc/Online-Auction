package com.Online_Auction.user_service.mapper;

import com.Online_Auction.user_service.domain.User;
import com.Online_Auction.user_service.dto.response.UserResponse;

public class UserMapper {
    public static UserResponse toUserResponse(User user) {
        UserResponse userResponse = new UserResponse();
        userResponse.setId(user.getId());
        userResponse.setBirthDate(user.getBirthDay());
        userResponse.setEmail(user.getEmail());
        userResponse.setEmailVerified(user.getEmailVerified());
        userResponse.setFullName(user.getFullName());
        userResponse.setIsSellerRequestSent(user.getIsSellerRequestSent());
        userResponse.setPassword(userResponse.getPassword());
        userResponse.setTotalNumberGoodReviews(user.getTotalNumberGoodReviews());
        userResponse.setTotalNumberReviews(userResponse.getTotalNumberReviews());
        userResponse.setUserRole(user.getRole());
        userResponse.setPassword(user.getPassword());
        return userResponse;
    }
}
