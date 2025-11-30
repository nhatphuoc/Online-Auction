package com.Online_Auction.auth_service.service;

import java.util.Map;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

@Service
public class RestTemplateNotificationService {

    private final RestTemplate restTemplate;

    @Value("${internal.key}")
    private String internalKey;

    private final String notificationBaseUrl = "http://localhost:8082/api/notify";

    public RestTemplateNotificationService(RestTemplate restTemplate) {
        this.restTemplate = restTemplate;
    }

    @Async
    public void sendEmail(String to, String subject, String body) {
        String url = notificationBaseUrl + "/email";

        Map<String, String> payload = Map.of(
                "to", to,
                "subject", subject,
                "body", body
        );

        HttpHeaders headers = new HttpHeaders();
        headers.set("X-Auth-Internal-Service", internalKey);

        HttpEntity<Map<String, String>> entity = new HttpEntity<>(payload, headers);

        try {
            restTemplate.postForObject(url, entity, Void.class);
        } catch (Exception ex) {
            throw new RuntimeException("Failed to send email via notification-service", ex);
        }
    }
}