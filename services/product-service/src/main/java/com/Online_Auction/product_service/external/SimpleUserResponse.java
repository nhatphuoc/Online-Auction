package com.Online_Auction.product_service.external;

import com.Online_Auction.product_service.config.security.UserPrincipal.UserRole;

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
