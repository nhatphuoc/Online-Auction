package com.Online_Auction.auth_service.external.client;

import java.util.Map;

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

import com.Online_Auction.auth_service.dto.request.RegisterUserRequest;
import com.Online_Auction.auth_service.dto.request.SignInRequest;
import com.Online_Auction.auth_service.external.response.ApiResponse;
import com.Online_Auction.auth_service.external.response.SimpleUserResponse;
import com.Online_Auction.auth_service.external.response.UserProfileResponse;
import com.fasterxml.jackson.databind.ObjectMapper;

@Service
public class UserServiceClient {

    private final RestTemplate restTemplate;

    @Value("${internal.key}")
    private String internalKey;

    private static final String BASE_URL = "http://localhost:8081/api/users";

    public UserServiceClient(RestTemplate restTemplate) {
        this.restTemplate = restTemplate;
    }

    /* ---------------------------------------------------------
     *  COMMON HEADER BUILDER
     * --------------------------------------------------------- */
    private HttpHeaders buildHeaders() {
        HttpHeaders headers = new HttpHeaders();
        headers.set("X-Auth-Internal-Service", internalKey);
        headers.setContentType(MediaType.APPLICATION_JSON);
        return headers;
    }

    /* ---------------------------------------------------------
     *  GENERIC CALLER — dùng chung cho mọi API
     * --------------------------------------------------------- */
    private <T> ApiResponse<T> callApi(
            String url,
            HttpMethod method,
            Object body,
            ParameterizedTypeReference<ApiResponse<T>> typeRef,
            Object... uriVars
    ) {
        try {
            HttpEntity<Object> entity = new HttpEntity<>(body, buildHeaders());

            ResponseEntity<ApiResponse<T>> response =
                    restTemplate.exchange(url, method, entity, typeRef, uriVars);

            return response.getBody();

        } catch (RestClientResponseException ex) {
            // ❗ Server trả về JSON error → parse ApiResponse từ response body
            try {
                ObjectMapper mapper = new ObjectMapper();
                return mapper.readValue(
                        ex.getResponseBodyAsString(),
                        mapper.getTypeFactory().constructParametricType(ApiResponse.class, Object.class)
                );
            } catch (Exception e) {
                return ApiResponse.fail("Failed to parse error response from user-service");
            }

        } catch (Exception ex) {
            return ApiResponse.fail("User-service unreachable: " + ex.getMessage());
        }
    }

    /* ---------------------------------------------------------
     *  API IMPLEMENTATION
     * --------------------------------------------------------- */

    /**
     * GET Simple User by Email
     */
    public ApiResponse<SimpleUserResponse> getUserByEmail(String email) {
        return callApi(
                BASE_URL + "/simple?email={email}",
                HttpMethod.GET,
                null,
                new ParameterizedTypeReference<>() {},
                email
        );
    }

    /**
     * Register User
     */
    public ApiResponse<Void> registerUser(RegisterUserRequest request) {
        return callApi(
                BASE_URL,
                HttpMethod.POST,
                request,
                new ParameterizedTypeReference<>() {}
        );
    }

    /**
     * Verify Email
     */
    public ApiResponse<Void> verifyEmail(String email) {
        return callApi(
                BASE_URL + "/verify-email",
                HttpMethod.POST,
                Map.of("email", email),
                new ParameterizedTypeReference<>() {}
        );
    }

    /**
     * Delete User by Email
     */
    public ApiResponse<Void> deleteUserByEmail(String email) {
        return callApi(
                BASE_URL,
                HttpMethod.DELETE,
                Map.of("email", email),
                new ParameterizedTypeReference<>() {}
        );
    }

    /**
     * Authenticate user
     */
    public ApiResponse<SimpleUserResponse> authenticateUser(SignInRequest request) {
        return callApi(
                BASE_URL + "/authenticate",
                HttpMethod.POST,
                request,
                new ParameterizedTypeReference<>() {}
        );
    }

    /**
     * Get logged-in user's profile
     */
    public ApiResponse<UserProfileResponse> getMyProfile() {
        return callApi(
                BASE_URL + "/profile/me",
                HttpMethod.GET,
                null,
                new ParameterizedTypeReference<>() {}
        );
    }

}