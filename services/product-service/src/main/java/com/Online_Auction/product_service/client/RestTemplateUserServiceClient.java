package com.Online_Auction.product_service.client;

import javax.management.RuntimeErrorException;

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
import com.fasterxml.jackson.databind.ObjectMapper;

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
        log.info("Calling user-service: GET {}", url);
        log.debug("Request header X-Auth-Internal-Service={}", internalKey);

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

            log.info("User-service response status={} for userId={}",
                    response.getStatusCode(), id);

            if (!response.getBody().isSuccess()) {
                throw new RuntimeException(
                        "User-service error: " + response.getBody().getMessage());
            }
            ObjectMapper mapper = new ObjectMapper();
            log.info("User: {}", mapper.writeValueAsString(response.getBody().getData()));
            return response.getBody().getData();
        } catch (HttpClientErrorException | HttpServerErrorException ex) {
            throw new RuntimeException(
                    "User-service error: " + ex.getStatusCode() + " - " + ex.getResponseBodyAsString(),
                    ex);
        } catch (Exception e) {
            throw new RuntimeException("Error");
        }
    }

}