package com.Online_Auction.product_service.mapper;

import org.springframework.stereotype.Component;

import com.Online_Auction.product_service.domain.Question;
import com.Online_Auction.product_service.dto.QuestionDTO;

import java.util.List;
import java.util.stream.Collectors;

@Component
public class QuestionMapper {

    private final AnswerMapper answerMapper;

    public QuestionMapper(AnswerMapper answerMapper) {
        this.answerMapper = answerMapper;
    }

    public QuestionDTO toDTO(Question question) {
        if (question == null) return null;

        return QuestionDTO.builder()
                .id(question.getId())
                .userId(question.getUserId())
                .content(question.getContent())
                .createdAt(question.getCreatedAt())
                .answer(answerMapper.toDTO(question.getAnswer()))
                .status(question.getStatus())
                .build();
    }

    public Question toEntity(QuestionDTO dto) {
        if (dto == null) return null;

        Question question = new Question();
        question.setId(dto.getId());
        question.setUserId(dto.getUserId());
        question.setContent(dto.getContent());
        question.setCreatedAt(dto.getCreatedAt());
        question.setAnswer(answerMapper.toEntity(dto.getAnswer()));
        question.setStatus(dto.getStatus());
        return question;
    }

    public List<QuestionDTO> toDTOList(List<Question> questions) {
        return questions.stream()
                .map(this::toDTO)
                .collect(Collectors.toList());
    }
}
