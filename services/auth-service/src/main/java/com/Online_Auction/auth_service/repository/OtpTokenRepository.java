package com.Online_Auction.auth_service.repository;

import java.util.Optional;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import com.Online_Auction.auth_service.domain.OtpToken;

@Repository
public interface OtpTokenRepository extends JpaRepository<OtpToken, Long>{
    Optional<OtpToken> findByEmail(String email);
    void deleteByEmail(String email);
}