package com.Online_Auction.product_service.repository;

import java.util.List;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import com.Online_Auction.product_service.domain.Question;

@Repository
public interface QuestionRepository extends JpaRepository<Question, Long> {

    List<Question> findByUserId(Long userId);
}
