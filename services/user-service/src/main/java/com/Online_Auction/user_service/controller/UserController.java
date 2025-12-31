package com.Online_Auction.user_service.controller;

import org.springframework.web.bind.annotation.RestController;

import com.Online_Auction.user_service.domain.UpgradeUser;
import com.Online_Auction.user_service.domain.User;
import com.Online_Auction.user_service.dto.request.RegisterUserRequest;
import com.Online_Auction.user_service.dto.request.SignInRequest;
import com.Online_Auction.user_service.dto.response.ApiResponse;
import com.Online_Auction.user_service.dto.response.SimpleUserResponse;
import com.Online_Auction.user_service.dto.response.StatusResponse;
import com.Online_Auction.user_service.dto.response.UserProfileResponse;
import com.Online_Auction.user_service.dto.response.UserSearchResponse;
import com.Online_Auction.user_service.mapper.UserMapper;
import com.Online_Auction.user_service.service.UserService;
import com.Online_Auction.user_service.service.UserUpgradeService;

import lombok.RequiredArgsConstructor;

import org.springframework.web.bind.annotation.RequestMapping;

import java.util.Map;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;

@RestController
@RequestMapping("/api/users")
@RequiredArgsConstructor
public class UserController {

    private final UserService userService;

    @GetMapping("/simple")
    public ResponseEntity<ApiResponse<SimpleUserResponse>> getSimpleUserByEmail(@RequestParam String email) {
        User user = userService.findByEmail(email);
        if (user == null)
            return ResponseEntity.badRequest().body(ApiResponse.fail("User not found"));

        return ResponseEntity.ok(
                ApiResponse.success(UserMapper.toSimpleUserResponse(user), "User fetched successfully"));
    }

    @GetMapping("/{id}/simple")
    public ResponseEntity<ApiResponse<SimpleUserResponse>> getSimpleUserById(@PathVariable Long id) {
        User user = userService.findById(id);
        if (user == null)
            return ResponseEntity.badRequest().body(ApiResponse.fail("User not found"));

        return ResponseEntity.ok(
                ApiResponse.success(UserMapper.toSimpleUserResponse(user), "User fetched successfully"));
    }

    @PostMapping
    public ResponseEntity<ApiResponse<Void>> registerUser(@RequestBody RegisterUserRequest entity) {
        boolean registerSuccess = userService.register(entity);

        if (!registerSuccess) {
            return ResponseEntity.badRequest()
                    .body(ApiResponse.fail("Fail to register user, email is already registered"));
        }

        return ResponseEntity.ok(ApiResponse.success(null, "Successfully registered user"));
    }

    @PostMapping("/verify-email")
    public ResponseEntity<ApiResponse<Void>> verifyEmail(@RequestBody Map<String, String> payload) {
        String email = payload.get("email");
        StatusResponse status = userService.verifyEmail(email);

        if (!status.isSuccess())
            return ResponseEntity.badRequest().body(ApiResponse.fail(status.getMessage()));

        return ResponseEntity.ok(ApiResponse.success(null, status.getMessage()));
    }

    @DeleteMapping
    public ResponseEntity<ApiResponse<Void>> deleteUserByEmail(@RequestBody Map<String, String> payload) {
        StatusResponse status = userService.deleteUserByEmail(payload.get("email"));

        if (!status.isSuccess())
            return ResponseEntity.badRequest().body(ApiResponse.fail(status.getMessage()));

        return ResponseEntity.ok(ApiResponse.success(null, status.getMessage()));
    }

    @PostMapping("/authenticate")
    public ResponseEntity<ApiResponse<SimpleUserResponse>> authenticateUser(@RequestBody SignInRequest request) {
        SimpleUserResponse response = userService.authenticateUser(request);

        return ResponseEntity.ok(ApiResponse.success(
                response,
                "Authentication successful"));
    }

    /**
     * GET /api/users/profile/me
     */
    @GetMapping("/profile/me")
    @PreAuthorize("hasAnyRole('BIDDER', 'SELLER', 'ADMIN')")
    public ResponseEntity<ApiResponse<UserProfileResponse>> getUserProfile() {
        User user = userService.getCurrentUser();

        if (user == null)
            return ResponseEntity.badRequest().body(ApiResponse.fail("User not authenticated"));

        return ResponseEntity.ok(
                ApiResponse.success(UserMapper.toUserProfileResponse(user), "Profile retrieved"));
    }

    @GetMapping("/search")
    @PreAuthorize("hasAnyRole('ADMIN', 'SELLER')")
    public Page<UserSearchResponse> searchUsers(
            @RequestParam(required = false) String keyword,
            @RequestParam(required = false) User.UserRole role,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size) {
        Pageable pageable = PageRequest.of(page, size, Sort.by("id").descending());
        return userService.searchUsers(keyword, role, pageable);
    }

    private final UserUpgradeService userUpgradeService;

    @PostMapping("/upgrade-to-seller")
    @PreAuthorize("hasAnyRole('BIDDER')")
    public ResponseEntity<?> requestUpgrade(@RequestParam String reason) {
        User user = userService.getCurrentUser();
        if (user == null)
            return ResponseEntity.badRequest().body(ApiResponse.fail("User not authenticated"));
        userUpgradeService.requestUpgradeToSeller(user.getId(), reason);
        return ResponseEntity.ok("Upgrade request submitted");
    }

    @PostMapping("/{requestId}/approve")
    @PreAuthorize("hasAnyRole('ADMIN')")
    public ResponseEntity<?> approve(@PathVariable Long requestId) {
        User admin = userService.getCurrentUser();

        if (admin == null)
            return ResponseEntity.badRequest().body(ApiResponse.fail("User not authenticated"));

        userUpgradeService.approveUpgrade(requestId, admin.getId());
        return ResponseEntity.ok("User upgraded to SELLER");
    }

    @PostMapping("/{requestId}/reject")
    @PreAuthorize("hasAnyRole('ADMIN')")
    public ResponseEntity<?> reject(@PathVariable Long requestId, @RequestParam(required = false) String rejectReason) {
        User admin = userService.getCurrentUser();

        if (admin == null)
            return ResponseEntity.badRequest().body(ApiResponse.fail("User not authenticated"));

        userUpgradeService.rejectUpgrade(requestId, admin.getId(), rejectReason);
        return ResponseEntity.ok("Upgrade request rejected");
    }

    @GetMapping("/upgrade-requests")
    @PreAuthorize("hasAnyRole('ADMIN')")
    public ResponseEntity<Page<UpgradeUser>> getPendingUpgradeRequests(
            @RequestParam(required = false) UpgradeUser.UpgradeStatus status,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size,
            @RequestParam(defaultValue = "createdAt") String sort,
            @RequestParam(defaultValue = "desc") String direction) {
        Sort sortOrder = Sort.by(
                direction.equalsIgnoreCase("desc") ? Sort.Direction.DESC : Sort.Direction.ASC,
                sort);

        Pageable pageable = PageRequest.of(page, size, sortOrder);

        Page<UpgradeUser> result;
        if (status != null) {
            result = userUpgradeService.getUpgradeRequestsByStatus(status, pageable);
        } else {
            result = userUpgradeService.getAllUpgradeRequests(pageable);
        }

        return ResponseEntity.ok(result);
    }

    @GetMapping
    public ResponseEntity<Page<UpgradeUser>> getAll(
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size,
            @RequestParam(defaultValue = "createdAt") String sort,
            @RequestParam(defaultValue = "desc") String direction) {
        Sort sortOrder = Sort.by(
                direction.equalsIgnoreCase("desc") ? Sort.Direction.DESC : Sort.Direction.ASC,
                sort);

        Pageable pageable = PageRequest.of(page, size, sortOrder);

        Page<UpgradeUser> result = userUpgradeService.getAllUpgradeRequests(pageable);

        return ResponseEntity.ok(result);
    }
}
