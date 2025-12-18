package com.Online_Auction.user_service.repository;

import java.util.Optional;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import com.Online_Auction.user_service.domain.User;

@Repository
public interface UserRepository extends JpaRepository<User, Long> {
    Optional<User> findByEmail(String email);

    void deleteByEmail(String email);

    @Query("""
                SELECT u FROM User u
                WHERE (:keyword IS NULL
                    OR LOWER(u.email) LIKE LOWER(CONCAT('%', :keyword, '%'))
                    OR LOWER(u.fullName) LIKE LOWER(CONCAT('%', :keyword, '%')))
                AND (:role IS NULL OR u.role = :role)
            """)
    Page<User> search(
            @Param("keyword") String keyword,
            @Param("role") User.UserRole role,
            Pageable pageable);
}