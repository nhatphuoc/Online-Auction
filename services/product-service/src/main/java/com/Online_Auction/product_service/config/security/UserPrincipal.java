package com.Online_Auction.product_service.config.security;

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
        BIDDER,
        SELLER,
        ADMIN
    }
}
