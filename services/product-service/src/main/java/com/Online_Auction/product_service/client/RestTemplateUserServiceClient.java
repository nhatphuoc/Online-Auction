package com.Online_Auction.product_service.client;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.client.HttpClientErrorException;
import org.springframework.web.client.HttpServerErrorException;
import org.springframework.web.client.RestTemplate;

import com.Online_Auction.product_service.dto.response.ApiResponse;
import com.Online_Auction.product_service.external.SimpleUserResponse;

import jakarta.annotation.PostConstruct;
import lombok.extern.slf4j.Slf4j;

@Service
@Slf4j
public class RestTemplateUserServiceClient {

    private final RestTemplate restTemplate;

    @Value("${internal.key}")
    private String internalKey;

    @Value("${USER_SERVICE_URL}")
    private String userServiceBaseUrl;

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

            ResponseEntity<ApiResponse<SimpleUserResponse>> response = restTemplate.exchange(
                    url,
                    HttpMethod.GET,
                    requestEntity,
                    new ParameterizedTypeReference<ApiResponse<SimpleUserResponse>>() {
                    });

            if (!response.getBody().isSuccess()) {
                throw new RuntimeException(
                        "User-service error: " + response.getBody().getMessage());
            }
            return response.getBody().getData();
        } catch (HttpClientErrorException | HttpServerErrorException ex) {
            throw new RuntimeException(
                    "User-service error: " + ex.getStatusCode() + " - " + ex.getResponseBodyAsString(),
                    ex);
        }
    }

}