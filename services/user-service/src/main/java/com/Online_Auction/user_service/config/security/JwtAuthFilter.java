package com.Online_Auction.user_service.config.security;

import java.io.IOException;
import java.util.List;

import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import com.Online_Auction.user_service.config.security.UserPrincipal.UserRole;

import io.jsonwebtoken.Claims;
import io.micrometer.common.util.StringUtils;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

@Component
public class JwtAuthFilter extends OncePerRequestFilter {

    private final InternalKeyValidator internalKeyValidator;
    private final TokenParser tokenParser;

    public JwtAuthFilter(InternalKeyValidator internalKeyValidator, TokenParser tokenParser) {
        this.internalKeyValidator = internalKeyValidator;
        this.tokenParser = tokenParser;
    }

    @Override
    protected void doFilterInternal(HttpServletRequest request,
            HttpServletResponse response,
            FilterChain filterChain)
            throws ServletException, IOException {

        // 1. Validate internal keys
        if (!internalKeyValidator.isValid(request)) {
            response.sendError(HttpServletResponse.SC_UNAUTHORIZED, "Invalid internal key");
            return;
        }

        // 2. Parse JWT từ header X-User-Token
        String token = request.getHeader("X-User-Token");
        if (!StringUtils.isBlank(token)) {
            Claims claims = tokenParser.parseClaims(token);

            // Lấy role từ claim
            String role = claims.get("role", String.class); // nếu role là List<String>

            // Chuyển sang GrantedAuthority
            List<SimpleGrantedAuthority> authorities = List.of(new SimpleGrantedAuthority(role));

            // Tạo principal
            UserPrincipal principal = new UserPrincipal(
                    Long.parseLong(claims.getSubject()),
                    claims.get("email", String.class),
                    UserRole.valueOf(role));

            // Tạo Authentication với authorities
            UsernamePasswordAuthenticationToken auth = new UsernamePasswordAuthenticationToken(
                    principal,
                    null,
                    authorities);

            // Set vào context
            SecurityContextHolder.getContext().setAuthentication(auth);
        } else {
            // Chuyển sang GrantedAuthority
            List<SimpleGrantedAuthority> authorities = List.of(new SimpleGrantedAuthority("ROLE_GATEWAY"));

            // Tạo Authentication với authorities
            UsernamePasswordAuthenticationToken auth = new UsernamePasswordAuthenticationToken(
                    null,
                    null,
                    authorities);

            // Set vào context
            SecurityContextHolder.getContext().setAuthentication(auth);
        }

        filterChain.doFilter(request, response);
    }
}
