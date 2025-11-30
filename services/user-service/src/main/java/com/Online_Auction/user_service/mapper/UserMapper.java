package com.Online_Auction.user_service.mapper;

import com.Online_Auction.user_service.domain.User;
import com.Online_Auction.user_service.dto.response.SimpleUserResponse;
import com.Online_Auction.user_service.dto.response.UserProfileResponse;

public class UserMapper {
    public static UserProfileResponse toUserProfileResponse(User user) {
        UserProfileResponse userResponse = new UserProfileResponse();
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

    public static SimpleUserResponse toSimpleUserResponse(User user) {
        SimpleUserResponse simpleUserResponse = new SimpleUserResponse();
        simpleUserResponse.setId(user.getId());
        simpleUserResponse.setEmail(user.getEmail());
        simpleUserResponse.setFullName(user.getFullName());
        simpleUserResponse.setUserRole(user.getRole());
        return simpleUserResponse;
    }
}
