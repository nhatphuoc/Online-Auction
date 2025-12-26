package com.Online_Auction.product_service.client;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.MediaType;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import com.Online_Auction.product_service.external.notification.EmailNotificationRequest;
import com.fasterxml.jackson.databind.ObjectMapper;

import lombok.extern.slf4j.Slf4j;

@Service
@Slf4j
public class RestTemplateNotificationServiceClient {

    private final RestTemplate restTemplate;

    @Value("${internal.key}")
    private String internalKey;

    @Value("${NOTIFICATION_SERVICE_URL}")
    private String notificationServiceBaseUrl;

    private ObjectMapper objectMapper = new ObjectMapper();

    public RestTemplateNotificationServiceClient(RestTemplate restTemplate) {
        this.restTemplate = restTemplate;
    }

    /**
     * Send email notification
     */
    public void sendEmail(EmailNotificationRequest request) {
        String url = notificationServiceBaseUrl + "/api/notify/email";

        try {
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.APPLICATION_JSON);
            headers.set("X-Auth-Internal-Service", internalKey);

            // 🔍 Log headers
            log.info("Request Headers: {}",
                    objectMapper.writeValueAsString(headers.toSingleValueMap()));

            // 🔍 Log body
            log.info("Request Body: {}",
                    objectMapper.writeValueAsString(request));

            HttpEntity<EmailNotificationRequest> requestEntity = new HttpEntity<>(request, headers);

            restTemplate.exchange(
                    url,
                    HttpMethod.POST,
                    requestEntity,
                    Void.class);
        } catch (Exception ex) {
            throw new RuntimeException("Failed to send email notification", ex);
        }
    }
}
