package com.Online_Auction.auth_service.config.filters;

import java.io.IOException;
import java.util.List;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.extern.slf4j.Slf4j;

@Component
@Slf4j
public class ApiGatewayAuthenticationFilter extends OncePerRequestFilter {

    private static final String HEADER_NAME = "X-Api-Gateway";

    @Value("${gateway.key}")
    private String apiGatewayKey;

    @Override
    protected void doFilterInternal(HttpServletRequest request,
            HttpServletResponse response,
            FilterChain filterChain) throws ServletException, IOException {

        String uri = request.getRequestURI();
        boolean isAuthEndpoint = uri.startsWith("/auth");

        if (!isAuthEndpoint) {
            String apiKey = request.getHeader(HEADER_NAME);
            if (apiKey == null || !apiKey.equals(apiGatewayKey)) {
                response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
                return;
            }
        }

        // Nếu không hợp lệ, SecurityConfig sẽ block access
        log.info("Incoming request");
        log.info("Method: {}", request.getMethod());
        log.info("URI: {}", request.getRequestURI());
        log.info("Authorization: {}", request.getHeader("X-Api-Gateway"));

        filterChain.doFilter(request, response);

        log.info("Response status: {}", response.getStatus());
    }
}