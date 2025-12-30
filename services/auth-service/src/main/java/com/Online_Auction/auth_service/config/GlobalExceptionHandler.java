package com.Online_Auction.auth_service.config;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;
import org.springframework.web.server.ResponseStatusException;

import com.Online_Auction.auth_service.external.response.ApiResponse;

import jakarta.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;

@RestControllerAdvice
@Slf4j
public class GlobalExceptionHandler {

        @ExceptionHandler(ResponseStatusException.class)
        public ResponseEntity<ApiResponse<Void>> handleResponseStatusException(
                        ResponseStatusException ex,
                        HttpServletRequest request) {

                log.warn("{} {} - {}",
                                request.getMethod(),
                                request.getRequestURI(),
                                ex.getReason());

                ApiResponse<Void> response = new ApiResponse<>();
                response.setSuccess(false);
                response.setMessage(ex.getReason());

                return ResponseEntity
                                .status(ex.getStatusCode())
                                .body(response);
        }
}
