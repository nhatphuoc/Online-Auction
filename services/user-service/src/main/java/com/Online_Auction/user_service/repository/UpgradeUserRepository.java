package com.Online_Auction.user_service.repository;

import java.util.Optional;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;

import com.Online_Auction.user_service.domain.UpgradeUser;

public interface UpgradeUserRepository extends JpaRepository<UpgradeUser, Long> {

    boolean existsByUserIdAndStatus(Long userId, UpgradeUser.UpgradeStatus status);

    Optional<UpgradeUser> findByIdAndStatus(Long id, UpgradeUser.UpgradeStatus status);

    Page<UpgradeUser> findAll(Pageable pageable);
}
