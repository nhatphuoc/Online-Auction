package com.Online_Auction.notification_service.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class SimpleUserResponse {
    private Long id;
    private String email;
    private String fullName;
    private UserRole userRole;

    public enum UserRole {
        ROLE_BIDDER,
        ROLE_MANAGER,
        ROLE_ADMIN
    }
}