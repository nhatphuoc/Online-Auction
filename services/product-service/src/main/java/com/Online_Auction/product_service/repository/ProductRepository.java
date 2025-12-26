package com.Online_Auction.product_service.repository;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;

import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.JpaSpecificationExecutor;
import org.springframework.data.jpa.repository.Lock;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import com.Online_Auction.product_service.domain.Product;

import jakarta.persistence.LockModeType;

@Repository
public interface ProductRepository extends JpaRepository<Product, Long>, JpaSpecificationExecutor<Product> {

    List<Product> findBySellerId(Long sellerId);

    List<Product> findByCategoryId(Long categoryId);

    @Lock(LockModeType.PESSIMISTIC_WRITE)
    @Query("SELECT p FROM Product p WHERE p.id = :id")
    Optional<Product> findByIdForUpdate(@Param("id") Long id);

    /* ================= HOMEPAGE ================= */

    // 1. Top 5 Auctions Gần Kết Thúc
    @Query("""
                SELECT p FROM Product p
                WHERE p.endAt > :now
                ORDER BY p.endAt ASC
            """)
    List<Product> findTop5EndingSoon(@Param("now") LocalDateTime now, Pageable pageable);

    // 2. Top 5 giá cao nhất
    @Query("""
                SELECT p FROM Product p
                WHERE p.endAt > :now
                ORDER BY p.currentPrice DESC
            """)
    List<Product> findTop5HighestPrice(@Param("now") LocalDateTime now, Pageable pageable);

    // 2. Top 5 nhiều lượt ra giá nhất — dùng @Query (join bids)
    @Query("""
                SELECT p FROM Product p
                WHERE p.endAt > :now
                ORDER BY p.bidCount DESC
            """)
    List<Product> findTop5MostBids(@Param("now") LocalDateTime now, Pageable pageable);

    /* ================= CATEGORY UPDATE ================= */

    @Modifying
    @Query("""
                UPDATE Product p
                SET p.categoryName = :categoryName,
                    p.parentCategoryId = :parentCategoryId,
                    p.parentCategoryName = :parentCategoryName
                WHERE p.categoryId = :categoryId
            """)
    int updateByCategoryId(
            @Param("categoryId") Long categoryId,
            @Param("categoryName") String categoryName,
            @Param("parentCategoryId") Long parentCategoryId,
            @Param("parentCategoryName") String parentCategoryName);

    /* ================= PARENT CATEGORY RENAME ================= */

    @Modifying
    @Query("""
                UPDATE Product p
                SET p.parentCategoryName = :parentCategoryName
                WHERE p.parentCategoryId = :parentCategoryId
            """)
    int updateParentCategoryName(
            @Param("parentCategoryId") Long parentCategoryId,
            @Param("parentCategoryName") String parentCategoryName);

    /* ================= EMAIL CRON-JOB ================= */
    @Query("""
                SELECT p FROM Product p
                WHERE p.endAt <= CURRENT_TIMESTAMP
                  AND (p.orderCreated = false
                  OR p.sentEmail = false)
            """)
    List<Product> findExpiredAuctionsForProcessing();

}
