package com.Online_Auction.auth_service.config.jwt;

import io.jsonwebtoken.*;
import io.jsonwebtoken.security.Keys;

import java.security.Key;
import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import com.Online_Auction.auth_service.external.response.UserResponse;

public class JwtUtils {
    private static final String SECRET_KEY = "X8f3N2p9V6yR1qT7Z4wM0bC5sH2kJ8lP";
    private static final long ACCESS_TOKEN_EXPIRATION = 24 * 60 * 60 * 1000; // 1 day
    private static final long REFRESH_TOKEN_EXPIRATION = 7 * 24 * 60 * 60 * 1000; // 7 day

    private static Key KEY = Keys.hmacShaKeyFor(SECRET_KEY.getBytes());

    // ==============================================
    // CLAIM HELPERS
    // ==============================================
    private static Claims getClaims(String token) {
        return Jwts.parserBuilder()
                .setSigningKey(KEY)
                .build()
                .parseClaimsJws(token)
                .getBody();
    }

    public static String getUsernameFromToken(String token) {
        return getClaims(token).getSubject();
    }

    public static List<String> getRole(String token) {
        return getClaims(token).get("role", List.class);
    }

    public static String getType(String token) {
        return getClaims(token).get("type", String.class);
    }

    // ==============================================
    // GENERATE TOKEN
    // ==============================================
    private static String generateToken(
        UserResponse user,
        Map<String, Object> claims,
        long expiration
    ) {
        return Jwts.builder()
                .setClaims(claims)
                .setSubject(String.valueOf(user.getId()))
                .setIssuedAt(new Date())
                .setExpiration(new Date(System.currentTimeMillis() + expiration))
                .signWith(KEY, SignatureAlgorithm.HS256)
                .compact();
    }

    public static String generateAccessToken(UserResponse user) {
        Map<String, Object> claims = new HashMap<>();
        claims.put("type", "access");
        claims.put("role", user.getUserRole());
        claims.put("email", user.getEmail());
        return generateToken(user, claims, ACCESS_TOKEN_EXPIRATION);
    }

    public static String generateRefreshToken(UserResponse user) {
        Map<String, Object> claims = new HashMap<>();
        claims.put("type", "refresh");
        return generateToken(user, claims, REFRESH_TOKEN_EXPIRATION);
    }

    // ==============================================
    // VALIDATE TOKEN
    // ==============================================
    public static boolean validateToken(String token) {
        try {
            Claims claims = getClaims(token); // if parsing fails, it will throw
            Date expiration = claims.getExpiration();
            String type = getType(token);
            return expiration.after(new Date()) && "access".equals(type);
        } catch (Exception e) {
            return false;
        }
    }
}
