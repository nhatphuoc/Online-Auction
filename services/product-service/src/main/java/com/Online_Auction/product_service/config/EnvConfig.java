package com.Online_Auction.product_service.config;

import org.springframework.context.annotation.Configuration;

import io.github.cdimascio.dotenv.Dotenv;

@Configuration
public class EnvConfig {

    static {
        // Load .env from project root: ../../.env
        Dotenv dotenv = Dotenv.configure()
                .directory("../../")   // adjust if needed
                .filename(".env")
                .load();

        dotenv.entries().forEach(entry ->
                System.setProperty(entry.getKey(), entry.getValue())
        );
    }
}