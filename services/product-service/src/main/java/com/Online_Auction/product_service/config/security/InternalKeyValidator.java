package com.Online_Auction.product_service.config.security;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import jakarta.servlet.http.HttpServletRequest;

@Component
public class InternalKeyValidator {

    @Value("${gateway.key}")
    private String apiGatewayKey;

    @Value("${internal.key}")
    private String internalServiceKey;

    public boolean isValid(HttpServletRequest request) {
        String gatewayKey = request.getHeader("X-Api-Gateway");
        String internalKey = request.getHeader("X-Auth-Internal-Service");

        return apiGatewayKey.equals(gatewayKey) || internalServiceKey.equals(internalKey);
    }
}
