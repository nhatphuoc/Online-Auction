package com.Online_Auction.product_service.repository;

import java.util.List;
import java.util.Optional;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.JpaSpecificationExecutor;
import org.springframework.data.jpa.repository.Lock;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import com.Online_Auction.product_service.domain.Product;
import com.Online_Auction.product_service.domain.Product.ProductStatus;

import jakarta.persistence.LockModeType;

@Repository
public interface ProductRepository extends JpaRepository<Product, Long>, JpaSpecificationExecutor<Product> {

    List<Product> findBySellerId(Long sellerId);

    List<Product> findByCategoryId(Long categoryId);

    @Query("SELECT p FROM Product p WHERE p.status = 'ACTIVE' AND p.endAt > CURRENT_TIMESTAMP")
    List<Product> findActiveProducts();

    @Lock(LockModeType.PESSIMISTIC_WRITE)
    @Query("SELECT p FROM Product p WHERE p.id = :id")
    Optional<Product> findByIdForUpdate(@Param("id") Long id);

    // 1. Top 5 gần kết thúc (status ACTIVE, current time < endAt)
    List<Product> findTop5ByStatusOrderByEndAtAsc(Product.ProductStatus status);

    // 3. Top 5 giá cao nhất
    List<Product> findTop5ByStatusOrderByCurrentPriceDesc(Product.ProductStatus status);

    // 2. Top 5 nhiều lượt ra giá nhất — dùng @Query (join bids)
    List<Product> findTop5ByStatusOrderByBidCountDesc(ProductStatus status);

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
}
