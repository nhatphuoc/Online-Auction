package com.Online_Auction.product_service.mapper;

import org.springframework.stereotype.Component;

import com.Online_Auction.product_service.domain.Favorite;
import com.Online_Auction.product_service.dto.FavoriteDTO;

import java.util.List;
import java.util.stream.Collectors;

@Component
public class FavoriteMapper {

    public FavoriteDTO toDTO(Favorite favorite) {
        if (favorite == null) return null;

        return FavoriteDTO.builder()
                .productId(favorite.getProductId())
                .userId(favorite.getUserId())
                .createdAt(favorite.getCreatedAt())
                .build();
    }

    public Favorite toEntity(FavoriteDTO dto) {
        if (dto == null) return null;

        Favorite favorite = new Favorite();
        favorite.setProductId(dto.getProductId());
        favorite.setUserId(dto.getUserId());
        favorite.setCreatedAt(dto.getCreatedAt());
        return favorite;
    }

    public List<FavoriteDTO> toDTOList(List<Favorite> favorites) {
        return favorites.stream()
                .map(this::toDTO)
                .collect(Collectors.toList());
    }
}
