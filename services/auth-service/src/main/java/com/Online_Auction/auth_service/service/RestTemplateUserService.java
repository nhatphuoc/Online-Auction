package com.Online_Auction.auth_service.service;

import java.util.Map;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import com.Online_Auction.auth_service.dto.request.RegisterRequest;
import com.Online_Auction.auth_service.external.response.StatusResponse;
import com.Online_Auction.auth_service.external.response.UserResponse;

@Service
public class RestTemplateUserService {

    private final RestTemplate restTemplate;

    @Value("${internal.key}")
    private String internalKey;

    private String userServiceUrl = "http://localhost:8081/api/users";

    public RestTemplateUserService(RestTemplate restTemplate) {
        this.restTemplate = restTemplate;
    }

    /**
     * Get user by email from user-service
     */
    public UserResponse getUserByEmail(String email) {
        String url = userServiceUrl + "?email={email}";

        try {
            HttpHeaders headers = new HttpHeaders();
            headers.set("X-Auth-Internal-Service", internalKey);

            HttpEntity<Void> requestEntity = new HttpEntity<>(headers);

            ResponseEntity<UserResponse> response = restTemplate.exchange(
                    url,
                    HttpMethod.GET,
                    requestEntity,
                    UserResponse.class,
                    email
            );

            return response.getBody();
        } catch (Exception ex) {
            throw new RuntimeException("Failed to get user from user-service", ex);
        }
    }

    /**
     * Register a new user via user-service
     */
    public StatusResponse registerUser(RegisterRequest request) {
        String url = userServiceUrl;

        try {
            HttpHeaders headers = new HttpHeaders();
            headers.set("X-Auth-Internal-Service", internalKey);
            headers.setContentType(MediaType.APPLICATION_JSON);

            HttpEntity<RegisterRequest> requestEntity = new HttpEntity<>(request, headers);

            ResponseEntity<StatusResponse> response = restTemplate.postForEntity(
                    url,
                    requestEntity,
                    StatusResponse.class
            );

            return response.getBody();
        } catch (Exception ex) {
            throw new RuntimeException("Failed to register user via user-service", ex);
        }
    }

    /**
     * Verify user email via user-service
     */
    public StatusResponse verifyEmail(String email) {
        String url = userServiceUrl + "/verify-email";

        try {
            HttpHeaders headers = new HttpHeaders();
            headers.set("X-Auth-Internal-Service", internalKey);
            headers.setContentType(MediaType.APPLICATION_JSON);

            Map<String, String> payload = Map.of("email", email);
            HttpEntity<Map<String, String>> requestEntity = new HttpEntity<>(payload, headers);

            ResponseEntity<StatusResponse> response = restTemplate.postForEntity(
                    url,
                    requestEntity,
                    StatusResponse.class
            );

            return response.getBody();
        } catch (Exception ex) {
            throw new RuntimeException("Failed to verify email via user-service", ex);
        }
    }
}
