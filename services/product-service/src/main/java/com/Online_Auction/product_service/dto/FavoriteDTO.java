package com.Online_Auction.product_service.dto;

import java.time.LocalDateTime;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter @Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class FavoriteDTO {

    private Long productId;
    private Long userId;
    private LocalDateTime createdAt;
}
