package com.online_auction.bidding_service.config.security;

import java.io.Serializable;

import lombok.AllArgsConstructor;
import lombok.Data;

@Data
@AllArgsConstructor
public class UserPrincipal implements Serializable {
    private Long userId;
    private String email;
    private UserRole role;
    
    public enum UserRole {
        // BIDDER, SELLER, ADMIN
        ROLE_BIDDER,
        ROLE_SELLER,
        ROLE_ADMIN
    }
}
