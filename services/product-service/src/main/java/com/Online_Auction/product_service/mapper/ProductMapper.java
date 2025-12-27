package com.Online_Auction.product_service.mapper;

import org.springframework.stereotype.Component;

import com.Online_Auction.product_service.domain.Product;
import com.Online_Auction.product_service.dto.response.ProductDTO;
import com.Online_Auction.product_service.dto.response.ProductListItemResponse;
import com.Online_Auction.product_service.dto.response.SimpleUserInfo;

import java.util.List;
import java.util.stream.Collectors;

@Component
public class ProductMapper {

    public ProductDTO toProductDTO(Product product, SimpleUserInfo seller, SimpleUserInfo highestBidder) {
        if (product == null)
            return null;

        return ProductDTO.builder()
                .id(product.getId())
                .name(product.getName())
                .thumbnailUrl(product.getThumbnailUrl())
                .images(product.getImages())
                .description(product.getDescription())
                .parentCategoryId(product.getParentCategoryId())
                .parentCategoryName(product.getParentCategoryName())
                .categoryId(product.getCategoryId())
                .categoryName(product.getCategoryName())
                .startingPrice(product.getStartingPrice())
                .currentPrice(product.getCurrentPrice())
                .buyNowPrice(product.getBuyNowPrice())
                .stepPrice(product.getStepPrice())
                .createdAt(product.getCreatedAt())
                .endAt(product.getEndAt())
                .autoExtend(product.isAutoExtend())
                .sellerInfo(seller)
                .highestBidder(highestBidder)
                .build();
    }

    public Product toEntity(ProductDTO dto) {
        if (dto == null)
            return null;

        Product product = new Product();
        product.setId(dto.getId());
        product.setName(dto.getName());
        product.setThumbnailUrl(dto.getThumbnailUrl());
        product.setImages(dto.getImages());
        product.setDescription(dto.getDescription());
        product.setCategoryId(dto.getCategoryId());
        product.setStartingPrice(dto.getStartingPrice());
        product.setCurrentPrice(dto.getCurrentPrice());
        product.setBuyNowPrice(dto.getBuyNowPrice());
        product.setStepPrice(dto.getStepPrice());
        product.setCreatedAt(dto.getCreatedAt());
        product.setEndAt(dto.getEndAt());
        product.setAutoExtend(dto.isAutoExtend());
        product.setSellerId(dto.getSellerId());
        return product;
    }

    public List<ProductDTO> toDTOList(List<Product> products, SimpleUserInfo seller, SimpleUserInfo highestBidder) {
        return products.stream()
                .map(p -> toProductDTO(p, seller, highestBidder))
                .collect(Collectors.toList());
    }

    public static ProductListItemResponse toListItem(Product p) {
        return new ProductListItemResponse(
                p.getId(),
                p.getThumbnailUrl(),
                p.getName(),
                p.getCurrentPrice(),
                p.getBuyNowPrice(),
                p.getCreatedAt(),
                p.getEndAt(),
                p.getBidCount(),

                // NEW FIELDS
                p.getParentCategoryId(),
                p.getParentCategoryName(),
                p.getCategoryId(),
                p.getCategoryName());
    }

}
