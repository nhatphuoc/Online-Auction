package com.Online_Auction.product_service.dto;

import com.Online_Auction.product_service.config.security.UserPrincipal.UserRole;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class SimpleUserInfo {
    private Long id;
    private String email;
    private String fullName;
    private UserRole userRole;
}
