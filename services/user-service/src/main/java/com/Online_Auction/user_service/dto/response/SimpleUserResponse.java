package com.Online_Auction.user_service.dto.response;

import com.Online_Auction.user_service.domain.User.UserRole;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class SimpleUserResponse {
    private long id;
    private String email;
    private String fullName;
    private UserRole userRole;
}
