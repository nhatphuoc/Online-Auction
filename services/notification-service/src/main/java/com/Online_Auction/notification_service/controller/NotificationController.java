package com.Online_Auction.notification_service.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RestController;

import com.Online_Auction.notification_service.dto.EmailRequest;
import com.Online_Auction.notification_service.service.EmailService;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;


@RestController
@RequestMapping("/api/notify")
public class NotificationController {
    
    @Autowired
    private EmailService emailService;

    @PostMapping("/email")
    public void sendEmail(@RequestBody EmailRequest emailRequest) {
        emailService.sendEmail(emailRequest);
    }
    
}
