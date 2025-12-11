package com.Online_Auction.product_service.specs;

import org.springframework.data.jpa.domain.Specification;

import com.Online_Auction.product_service.domain.Product;

public class ProductSpecs {

    public static Specification<Product> hasParentCategory(Long parentCategoryId) {
        return (root, query, cb) ->
                parentCategoryId == null ? null :
                        cb.equal(root.get("parentCategoryId"), parentCategoryId);
    }

    public static Specification<Product> hasCategory(Long categoryId) {
        return (root, query, cb) ->
                categoryId == null ? null :
                        cb.equal(root.get("categoryId"), categoryId);
    }

    public static Specification<Product> hasNamePrefix(String queryStr) {
        return (root, query, cb) ->
                (queryStr == null || queryStr.isBlank()) ? null :
                        cb.like(cb.lower(root.get("name")), queryStr.toLowerCase() + "%");
    }

    public static Specification<Product> isActive() {
        return (root, query, cb) ->
                cb.equal(root.get("status"), Product.ProductStatus.ACTIVE);
    }
}
