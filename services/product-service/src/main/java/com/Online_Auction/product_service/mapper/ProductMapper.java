package com.Online_Auction.product_service.mapper;

import org.springframework.stereotype.Component;

import com.Online_Auction.product_service.domain.Product;
import com.Online_Auction.product_service.domain.ProductStatus;
import com.Online_Auction.product_service.dto.ProductDTO;
import com.Online_Auction.product_service.dto.QuestionDTO;
import com.Online_Auction.product_service.dto.SimpleUserInfo;

import java.util.List;
import java.util.stream.Collectors;

@Component
public class ProductMapper {

    private final QuestionMapper questionMapper;

    public ProductMapper(QuestionMapper questionMapper) {
        this.questionMapper = questionMapper;
    }

    public ProductDTO toProductDTO(Product product, SimpleUserInfo seller, SimpleUserInfo highestBidder) {
        if (product == null) return null;

        List<QuestionDTO> questions = product.getQuestions() == null ? List.of()
                : product.getQuestions().stream()
                .map(questionMapper::toDTO)
                .collect(Collectors.toList());

        return ProductDTO.builder()
                .id(product.getId())
                .name(product.getName())
                .thumbnailUrl(product.getThumbnailUrl())
                .images(product.getImages())
                .description(product.getDescription())
                .categoryId(product.getCategoryId())
                .startingPrice(product.getStartingPrice())
                .currentPrice(product.getCurrentPrice())
                .buyNowPrice(product.getBuyNowPrice())
                .stepPrice(product.getStepPrice())
                .createdAt(product.getCreatedAt())
                .endAt(product.getEndAt())
                .autoExtend(product.isAutoExtend())
                .sellerInfo(seller)
                .highestBidder(highestBidder)
                .questions(questions)
                .status(product.getStatus() != null ? product.getStatus() : ProductStatus.ACTIVE)
                .build();
    }

    public Product toEntity(ProductDTO dto) {
        if (dto == null) return null;

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
        product.setQuestions(dto.getQuestions() == null ? null :
                dto.getQuestions().stream()
                        .map(questionMapper::toEntity)
                        .collect(Collectors.toList()));
        product.setStatus(dto.getStatus());
        return product;
    }

    public List<ProductDTO> toDTOList(List<Product> products, SimpleUserInfo seller, SimpleUserInfo highestBidder) {
        return products.stream()
                .map(p -> toProductDTO(p, seller, highestBidder))
                .collect(Collectors.toList());
    }
}
