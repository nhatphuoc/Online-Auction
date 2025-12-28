package com.Online_Auction.notification_service.client;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClientResponseException;
import org.springframework.web.client.RestTemplate;

import com.Online_Auction.notification_service.dto.SimpleUserResponse;

import tools.jackson.databind.ObjectMapper;

@Service
public class UserServiceClient {

    private final RestTemplate restTemplate;

    @Value("${internal.key}")
    private String internalKey;

    @Value("${USER_SERVICE_URL}")
    private String userServiceUrl; // Injected from application.yaml

    public UserServiceClient(RestTemplate restTemplate) {
        this.restTemplate = restTemplate;
    }

    /*
     * ---------------------------------------------------------
     * COMMON HEADER BUILDER
     * ---------------------------------------------------------
     */
    private HttpHeaders buildHeaders() {
        HttpHeaders headers = new HttpHeaders();
        headers.set("X-Auth-Internal-Service", internalKey);
        headers.setContentType(MediaType.APPLICATION_JSON);
        return headers;
    }

    /*
     * ---------------------------------------------------------
     * GENERIC CALLER — dùng chung cho mọi API
     * ---------------------------------------------------------
     */
    private <T> ApiResponse<T> callApi(
            String url,
            HttpMethod method,
            Object body,
            ParameterizedTypeReference<ApiResponse<T>> typeRef,
            Object... uriVars) {
        try {
            HttpEntity<Object> entity = new HttpEntity<>(body, buildHeaders());

            ResponseEntity<ApiResponse<T>> response = restTemplate.exchange(url, method, entity, typeRef, uriVars);

            return response.getBody();

        } catch (RestClientResponseException ex) {
            // Server trả về JSON error → parse ApiResponse từ response body
            try {
                ObjectMapper mapper = new ObjectMapper();
                return mapper.readValue(
                        ex.getResponseBodyAsString(),
                        mapper.getTypeFactory().constructParametricType(ApiResponse.class, Object.class));
            } catch (Exception e) {
                return ApiResponse.fail("Failed to parse error response from user-service");
            }

        } catch (Exception ex) {
            return ApiResponse.fail("User-service unreachable: " + ex.getMessage());
        }
    }

    /*
     * ---------------------------------------------------------
     * API IMPLEMENTATION
     * ---------------------------------------------------------
     */

    /**
     * GET Simple User by ID
     */
    public ApiResponse<SimpleUserResponse> getUserById(Long id) {
        return callApi(
                userServiceUrl + "/{id}/simple", // Use injected URL
                HttpMethod.GET,
                null,
                new ParameterizedTypeReference<>() {
                },
                id);
    }
}
