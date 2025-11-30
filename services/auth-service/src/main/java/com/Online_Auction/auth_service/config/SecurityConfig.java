package com.Online_Auction.auth_service.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter;

import com.Online_Auction.auth_service.config.filters.ApiGatewayAuthenticationFilter;

@Configuration
public class SecurityConfig {

    private final ApiGatewayAuthenticationFilter apiGatewayFilter;

    public SecurityConfig(ApiGatewayAuthenticationFilter apiGatewayFilter) {
        this.apiGatewayFilter = apiGatewayFilter;
    }

    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        http
            .csrf(c -> c.disable())
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/auth/**").permitAll()
                .anyRequest().denyAll()
            );

        http.addFilterBefore(apiGatewayFilter, UsernamePasswordAuthenticationFilter.class);
        
        return http.build();
    }
}