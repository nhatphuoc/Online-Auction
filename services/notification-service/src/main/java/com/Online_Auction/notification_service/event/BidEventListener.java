package com.Online_Auction.notification_service.event;

import java.util.Objects;

import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.stereotype.Service;

import com.Online_Auction.notification_service.client.ApiResponse;
import com.Online_Auction.notification_service.client.UserServiceClient;
import com.Online_Auction.notification_service.config.RabbitMQConfig;
import com.Online_Auction.notification_service.dto.EmailRequest;
import com.Online_Auction.notification_service.dto.SimpleUserResponse;
import com.Online_Auction.notification_service.service.*;
import lombok.RequiredArgsConstructor;

@Service
@RequiredArgsConstructor
public class BidEventListener {

    private final EmailService emailService;
    private final UserServiceClient userServiceClient;

    @RabbitListener(queues = RabbitMQConfig.BID_SUCCESS_QUEUE)
    public void handleBidSuccess(BidPlacedEvent event) {
        System.out.println("Event: " + event);
        // push websocket notifications

        // send emails to the current
        if (Objects.nonNull(event.getBidderId())) {
            ApiResponse<SimpleUserResponse> response = userServiceClient.getUserById(event.getBidderId());
            System.out.println(response);
            if (response.isSuccess()) {
                SimpleUserResponse data = response.getData();
                EmailRequest emailRequest = new EmailRequest();
                emailRequest.setTo(data.getEmail());
                emailRequest.setSubject("Bidding Notification");
                emailRequest.setBody("Bidding successfully");
                emailService.sendEmail(emailRequest);
            }
        }

        // send emails to the previous
        if (Objects.nonNull(event.getPreviousHighestBidder())) {
            ApiResponse<SimpleUserResponse> response = userServiceClient.getUserById(event.getPreviousHighestBidder());
            if (response.isSuccess()) {
                SimpleUserResponse data = response.getData();
                EmailRequest emailRequest = new EmailRequest();
                emailRequest.setTo(data.getEmail());
                emailRequest.setSubject("Bidding Notification");
                emailRequest.setBody("Bidding successfully");
                emailService.sendEmail(emailRequest);
            }
        }   
    }
}