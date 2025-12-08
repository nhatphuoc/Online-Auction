package com.Online_Auction.auth_service.external.response;



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
