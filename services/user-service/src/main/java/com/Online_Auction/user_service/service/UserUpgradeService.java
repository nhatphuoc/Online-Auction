package com.Online_Auction.user_service.service;

import java.time.LocalDateTime;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;

import com.Online_Auction.user_service.domain.UpgradeUser;
import com.Online_Auction.user_service.domain.User;
import com.Online_Auction.user_service.repository.UpgradeUserRepository;
import com.Online_Auction.user_service.repository.UserRepository;

import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;

@Service
@RequiredArgsConstructor
public class UserUpgradeService {

    private final UserRepository userRepository;
    private final UpgradeUserRepository upgradeUserRepository;

    @Transactional
    public void requestUpgradeToSeller(Long userId, String reason) {
        User user = userRepository.findById(userId)
                .orElseThrow(() -> new RuntimeException("User not found"));

        if (user.getRole() != User.UserRole.ROLE_BIDDER) {
            throw new RuntimeException("Only BIDDER can request upgrade");
        }

        if (upgradeUserRepository.existsByUserIdAndStatus(
                userId, UpgradeUser.UpgradeStatus.PENDING)) {
            throw new RuntimeException("Upgrade request already pending");
        }

        UpgradeUser request = UpgradeUser.builder()
                .user(user)
                .status(UpgradeUser.UpgradeStatus.PENDING)
                .reason(reason)
                .createdAt(LocalDateTime.now())
                .build();

        user.setIsSellerRequestSent(true);

        upgradeUserRepository.save(request);
        userRepository.save(user);
    }

    @Transactional
    public void approveUpgrade(Long requestId, Long adminId) {
        UpgradeUser request = upgradeUserRepository
                .findByIdAndStatus(requestId, UpgradeUser.UpgradeStatus.PENDING)
                .orElseThrow(() -> new RuntimeException("Request not found"));

        User user = request.getUser();
        user.setRole(User.UserRole.ROLE_SELLER);
        user.setIsSellerRequestSent(false);

        request.setStatus(UpgradeUser.UpgradeStatus.APPROVED);
        request.setReviewedAt(LocalDateTime.now());
        request.setReviewedByAdminId(adminId);

        userRepository.save(user);
        upgradeUserRepository.save(request);
    }

    public Page<UpgradeUser> getAllUpgradeRequests(Pageable pageable) {
        return upgradeUserRepository.findAll(pageable);
    }
}
