package com.Online_Auction.user_service.domain;

import java.time.LocalDateTime;

import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.FetchType;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.Table;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Entity
@Table(name = "user_upgrade_requests")
@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class UpgradeUser {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "user_id", nullable = false)
    private User user;

    @Enumerated(EnumType.STRING)
    private UpgradeStatus status;

    private String reason;

    private LocalDateTime createdAt;

    private LocalDateTime reviewedAt;

    private Long reviewedByAdminId;

    public enum UpgradeStatus {
        PENDING,
        APPROVED,
        REJECTED
    }
}
