package com.online_auction.bidding_service.config.security;

import java.io.IOException;
import java.util.List;

import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import com.online_auction.bidding_service.config.security.UserPrincipal.UserRole;

import io.jsonwebtoken.Claims;
import io.micrometer.common.util.StringUtils;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.extern.slf4j.Slf4j;

@Component
@Slf4j
public class JwtAuthFilter extends OncePerRequestFilter {

    private final InternalKeyValidator internalKeyValidator;
    private final TokenParser tokenParser;

    public JwtAuthFilter(InternalKeyValidator internalKeyValidator, TokenParser tokenParser) {
        this.internalKeyValidator = internalKeyValidator;
        this.tokenParser = tokenParser;
    }

    @Override
    protected boolean shouldNotFilter(HttpServletRequest request) {
        String path = request.getServletPath();

        return path.startsWith("/swagger-ui")
                || path.startsWith("/v3/api-docs")
                || path.equals("/swagger-ui.html");
    }

    @Override
    protected void doFilterInternal(HttpServletRequest request,
            HttpServletResponse response,
            FilterChain filterChain)
            throws ServletException, IOException {

        String method = request.getMethod();
        String uri = request.getRequestURI();
        String query = request.getQueryString();

        log.info("➡️ Incoming request: {} {}{}",
                method,
                uri,
                query != null ? "?" + query : "");

        // 1. Validate internal keys
        if (!internalKeyValidator.isValid(request)) {
            response.sendError(HttpServletResponse.SC_UNAUTHORIZED, "Invalid internal key");
            return;
        }

        // 2. Parse JWT từ header X-User-Token
        String token = request.getHeader("X-User-Token");
        System.out.println("Token: " + token);
        if (!StringUtils.isBlank(token)) {
            Claims claims = tokenParser.parseClaims(token);

            // Lấy role từ claim
            String role = claims.get("role", String.class); // nếu role là List<String>
            System.out.println("Role: " + role);

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
