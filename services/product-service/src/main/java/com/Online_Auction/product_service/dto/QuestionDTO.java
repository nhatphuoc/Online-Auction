package com.Online_Auction.product_service.dto;

import java.time.LocalDateTime;

import com.Online_Auction.product_service.domain.QuestionStatus;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter @Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class QuestionDTO {

    private Long id;

    private Long userId;

    private String content;

    private LocalDateTime createdAt;

    private AnswerDTO answer;

    private QuestionStatus status;
}
