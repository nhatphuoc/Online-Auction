package com.Online_Auction.user_service.config.security;

import java.io.IOException;

import org.springframework.web.filter.OncePerRequestFilter;

import com.Online_Auction.user_service.config.InternalConfig;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

public class InternalAuthFilter extends OncePerRequestFilter {

    private final InternalConfig internalConfig;

    public InternalAuthFilter(InternalConfig internalConfig) {
        this.internalConfig = internalConfig;
    }

    @Override
    protected void doFilterInternal(HttpServletRequest request,
                                    HttpServletResponse response,
                                    FilterChain filterChain)
            throws ServletException, IOException {

        String header = request.getHeader("X-Auth-Internal-Service");

        if (header == null || !header.equals(internalConfig.getKey())) {
            response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
            response.getWriter().write("Unauthorized: invalid internal service key");
            return;
        }

        filterChain.doFilter(request, response);
    }
}