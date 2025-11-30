package com.Online_Auction.product_service.service;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import com.Online_Auction.product_service.domain.Favorite;
import com.Online_Auction.product_service.dto.FavoriteDTO;
import com.Online_Auction.product_service.mapper.FavoriteMapper;
import com.Online_Auction.product_service.repository.FavoriteRepository;

import java.time.LocalDateTime;
import java.util.List;

@Service
@RequiredArgsConstructor
public class FavoriteService {

    private final FavoriteRepository favoriteRepository;
    private final FavoriteMapper favoriteMapper;

    @Transactional
    public void addFavorite(Long userId, Long productId) {
        if (favoriteRepository.existsByUserIdAndProductId(userId, productId)) 
            return;

        Favorite favorite = Favorite.builder()
                .userId(userId)
                .productId(productId)
                .createdAt(LocalDateTime.now())
                .build();

        favoriteRepository.save(favorite);
    }

    @Transactional
    public void removeFavorite(Long userId, Long productId) {
        favoriteRepository.deleteByUserIdAndProductId(userId, productId);
    }

    @Transactional(readOnly = true)
    public List<FavoriteDTO> listFavorites(Long userId) {
        return favoriteMapper.toDTOList(favoriteRepository.findByUserId(userId));
    }
}
