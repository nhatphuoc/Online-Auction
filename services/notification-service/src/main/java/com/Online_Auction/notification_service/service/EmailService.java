package com.Online_Auction.notification_service.service;

import com.Online_Auction.notification_service.dto.EmailRequest;

public interface EmailService {
    void sendEmail(EmailRequest request);
}