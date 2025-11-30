package com.Online_Auction.product_service.mapper;

import org.springframework.stereotype.Component;

import com.Online_Auction.product_service.domain.Answer;
import com.Online_Auction.product_service.dto.AnswerDTO;

@Component
public class AnswerMapper {

    public AnswerDTO toDTO(Answer answer) {
        if (answer == null) return null;

        return AnswerDTO.builder()
                .id(answer.getId())
                .sellerId(answer.getSellerId())
                .message(answer.getMessage())
                .createdAt(answer.getCreatedAt())
                .build();
    }

    public Answer toEntity(AnswerDTO dto) {
        if (dto == null) return null;

        Answer answer = new Answer();
        answer.setId(dto.getId());
        answer.setSellerId(dto.getSellerId());
        answer.setMessage(dto.getMessage());
        answer.setCreatedAt(dto.getCreatedAt());
        return answer;
    }
}
