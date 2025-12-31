package com.Online_Auction.product_service.client;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import com.Online_Auction.product_service.external.order.CreateOrderRequest;
import com.Online_Auction.product_service.external.order.OrderResponse;
import com.fasterxml.jackson.databind.ObjectMapper;

import lombok.extern.slf4j.Slf4j;

@Service
@Slf4j
public class RestTemplateOrderServiceClient {

    private final RestTemplate restTemplate;

    @Value("${internal.key}")
    private String internalKey;

    @Value("${ORDER_SERVICE_URL}")
    private String orderServiceBaseUrl;

    public RestTemplateOrderServiceClient(RestTemplate restTemplate) {
        this.restTemplate = restTemplate;
    }

    /**
     * Create order after auction finished
     */
    public OrderResponse createOrder(CreateOrderRequest request) {
        String url = orderServiceBaseUrl + "/order";

        ObjectMapper objectMapper = new ObjectMapper();

        try {
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.APPLICATION_JSON);
            headers.set("X-Auth-Internal-Service", internalKey);

            HttpEntity<CreateOrderRequest> requestEntity = new HttpEntity<>(request, headers);

            // 🔍 Log headers
            log.info("Request Headers: {}",
                    objectMapper.writeValueAsString(headers.toSingleValueMap()));

            // 🔍 Log body
            log.info("Request Body: {}",
                    objectMapper.writeValueAsString(request));

            ResponseEntity<OrderResponse> response = restTemplate.exchange(
                    url,
                    HttpMethod.POST,
                    requestEntity,
                    OrderResponse.class);

            // 🔍 Log response body
            log.info("Response Body: {}",
                    objectMapper.writeValueAsString(response.getBody()));

            return response.getBody();

        } catch (Exception ex) {
            log.error("Failed to create order", ex);
            throw new RuntimeException("Failed to create order", ex);
        }
    }
}
