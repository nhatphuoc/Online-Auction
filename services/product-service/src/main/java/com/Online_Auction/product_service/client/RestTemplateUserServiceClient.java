package com.Online_Auction.product_service.client;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import com.Online_Auction.product_service.external.SimpleUserResponse;

@Service
public class RestTemplateUserServiceClient {

    private final RestTemplate restTemplate;

    @Value("${internal.key}")
    private String internalKey;

    private final String userServiceBaseUrl = "http://localhost:8081/api/users";

    public RestTemplateUserServiceClient(RestTemplate restTemplate) {
        this.restTemplate = restTemplate;
    }

    /**
     * Get SimpleUserResponse by id from user-service
     */
    public SimpleUserResponse getUserById(long id) {
        String url = userServiceBaseUrl + String.format("/%s/simple", id);

        try {
            HttpHeaders headers = new HttpHeaders();
            headers.set("X-Auth-Internal-Service", internalKey);

            HttpEntity<Void> requestEntity = new HttpEntity<>(headers);

            ResponseEntity<SimpleUserResponse> response = restTemplate.exchange(
                    url,
                    HttpMethod.GET,
                    requestEntity,
                    SimpleUserResponse.class
            );

            return response.getBody();
        } catch (Exception ex) {
            throw new RuntimeException("Failed to get user from user-service", ex);
        }
    }

}