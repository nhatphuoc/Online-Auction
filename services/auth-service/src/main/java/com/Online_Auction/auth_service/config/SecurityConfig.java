package com.Online_Auction.auth_service.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter;
import org.springframework.web.cors.CorsConfiguration;
import org.springframework.web.cors.CorsConfigurationSource;
import org.springframework.web.cors.UrlBasedCorsConfigurationSource;

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
                .cors(Customizer.withDefaults())
                .authorizeHttpRequests(auth -> auth
                        // 1. Auth endpoints
                        .requestMatchers(HttpMethod.OPTIONS, "/auth/**").permitAll()

                        // 2. Swagger / OpenAPI
                        .requestMatchers(
                                "/swagger-ui/**",
                                "/v3/api-docs/**",
                                "/swagger-ui.html")
                        .permitAll()

                        // 3. Everything else must pass the gateway filter
                        .anyRequest().authenticated());

        http.addFilterBefore(apiGatewayFilter, UsernamePasswordAuthenticationFilter.class);

        return http.build();
    }

    @Bean
    public CorsConfigurationSource corsConfigurationSource() {
        CorsConfiguration config = new CorsConfiguration();
        config.setAllowCredentials(true);

        // Allowed origins (your frontend)
        config.addAllowedOriginPattern("http://localhost:*");

        // Allowed headers
        config.addAllowedHeader("*");

        // Allowed methods
        config.addAllowedMethod("*");

        UrlBasedCorsConfigurationSource source = new UrlBasedCorsConfigurationSource();
        source.registerCorsConfiguration("/**", config);
        return source;
    }

}