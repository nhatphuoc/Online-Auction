package com.Online_Auction.user_service.mapper;

import com.Online_Auction.user_service.domain.User;
import com.Online_Auction.user_service.dto.response.SimpleUserResponse;
import com.Online_Auction.user_service.dto.response.UserProfileResponse;

public class UserMapper {
    
    public static UserProfileResponse toUserProfileResponse(User user) {
        if (user == null) return null;

        return new UserProfileResponse(
                user.getId(),
                user.getFullName(),
                user.getEmail(),
                user.getBirthDay(),
                user.getRole(),
                user.getTotalNumberReviews(),
                user.getTotalNumberGoodReviews(),
                user.getEmailVerified(),
                user.getIsSellerRequestSent()
        );
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
