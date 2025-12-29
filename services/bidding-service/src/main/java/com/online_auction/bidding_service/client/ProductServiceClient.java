package com.online_auction.bidding_service.client;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import com.online_auction.bidding_service.dto.request.ProductBidRequest;
import com.online_auction.bidding_service.dto.response.ApiResponse;
import com.online_auction.bidding_service.dto.response.ProductBidSuccessData;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Service
@RequiredArgsConstructor
@Slf4j
public class ProductServiceClient {

    private final RestTemplate restTemplate;

    @Value("${PRODUCT_SERVICE_URL}")
    private String productServiceUrl; // e.g., http://localhost:8085/api/products

    @Value("${internal.key}")
    private String internalKey;

    /**
     * Gửi yêu cầu đặt giá tới Product-Service.
     * Luôn forward X-User-Token sang cho product-service.
     */
    public ApiResponse<ProductBidSuccessData> placeBidToProductService(
            Long productId,
            ProductBidRequest request,
            String userJwt) {

        String url = productServiceUrl + "/" + productId + "/bids";

        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);

        // Forward JWT to product-service
        headers.set("X-User-Token", userJwt);
        headers.set("X-Auth-Internal-Service", internalKey);

        HttpEntity<ProductBidRequest> entity = new HttpEntity<>(request, headers);

        try {
            log.info("Sending POST request to {}", url);
            log.info("Request headers: {}", entity.getHeaders());
            log.info("Request body: {}", entity.getBody());

            ResponseEntity<ApiResponse<ProductBidSuccessData>> response = restTemplate.exchange(
                    url,
                    HttpMethod.POST,
                    entity,
                    new ParameterizedTypeReference<>() {
                    });

            log.info("Response status: {}", response.getStatusCode());
            log.info("Response headers: {}", response.getHeaders());
            log.info("Response body: {}", response.getBody());

            return response.getBody(); // Có thể null → BiddingService xử lý tiếp
        } catch (Exception ex) {
            // Trả về lỗi theo format ApiResponse
            return ApiResponse.fail("PRODUCT_SERVICE_ERROR: " + ex.getMessage());
        }
    }
}
