package com.Online_Auction.product_service.service;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import com.Online_Auction.product_service.domain.Answer;
import com.Online_Auction.product_service.domain.Question;
import com.Online_Auction.product_service.dto.QuestionDTO;
import com.Online_Auction.product_service.dto.request.AnswerCreateRequest;
import com.Online_Auction.product_service.dto.request.QuestionCreateRequest;
import com.Online_Auction.product_service.mapper.QuestionMapper;
import com.Online_Auction.product_service.repository.ProductRepository;
import com.Online_Auction.product_service.repository.QuestionRepository;

import java.time.LocalDateTime;

@Service
@RequiredArgsConstructor
public class QuestionService {

    private final ProductRepository productRepository;
    private final QuestionRepository questionRepository;
    private final QuestionMapper questionMapper;

    @Transactional
    public QuestionDTO askQuestion(Long userId, Long productId, QuestionCreateRequest request) {
        var product = productRepository.findById(productId)
                .orElseThrow(() -> new IllegalArgumentException("Product not found"));

        Question question = Question.builder()
                .userId(userId)
                .content(request.getContent())
                .createdAt(LocalDateTime.now())
                .build();

        product.getQuestions().add(question);
        productRepository.save(product);

        return questionMapper.toDTO(question);
    }

    @Transactional
    public QuestionDTO answerQuestion(Long sellerId, Long questionId, AnswerCreateRequest request) {
        Question question = questionRepository.findById(questionId)
                .orElseThrow(() -> new IllegalArgumentException("Question not found"));

        // TODO: check sellerId match product.sellerId
        Answer answer = Answer.builder()
                .sellerId(sellerId)
                .message(request.getMessage())
                .createdAt(LocalDateTime.now())
                .build();

        question.setAnswer(answer);
        questionRepository.save(question);

        return questionMapper.toDTO(question);
    }
}
