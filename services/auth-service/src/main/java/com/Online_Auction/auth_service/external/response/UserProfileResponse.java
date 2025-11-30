package com.Online_Auction.auth_service.external.response;

import java.time.LocalDate;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class UserProfileResponse {
    
    private Long id;
    private String password;
    private String email;
    private String fullName;
    private LocalDate birthDate;
    private int totalNumberReviews = 0;
    private int totalNumberGoodReviews = 0;
    private Boolean emailVerified = false;
    private Boolean isSellerRequestSent = false;
    private UserRole userRole;

    public enum UserRole {
        BIDDER,
        SELLER,
        ADMIN
    }
}
