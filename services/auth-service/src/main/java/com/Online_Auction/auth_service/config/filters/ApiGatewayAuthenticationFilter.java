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

@Component
public class ApiGatewayAuthenticationFilter extends OncePerRequestFilter {

    private static final String HEADER_NAME = "X-Api-Gateway";

    @Value("${gateway.key}")
    private String apiGatewayKey;

    @Override
    protected void doFilterInternal(HttpServletRequest request,
            HttpServletResponse response,
            FilterChain filterChain) throws ServletException, IOException {

        String apiKey = request.getHeader(HEADER_NAME);

        if (apiKey != null && apiKey.equals(apiGatewayKey)) {
            // Tạo Authentication đơn giản với role GATEWAY
            Authentication authentication = new UsernamePasswordAuthenticationToken(
                    "API_GATEWAY",
                    null,
                    List.of(new SimpleGrantedAuthority("ROLE_GATEWAY")));
            SecurityContextHolder.getContext().setAuthentication(authentication);
        }

        // Nếu không hợp lệ, SecurityConfig sẽ block access
        filterChain.doFilter(request, response);
    }
}