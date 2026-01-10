package com.online_auction.bidding_service.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import com.online_auction.bidding_service.domain.User;

@Repository
public interface UserRepository extends JpaRepository<User, Long> {

}