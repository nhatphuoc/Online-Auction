package com.Online_Auction.auth_service;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.scheduling.annotation.EnableAsync;

import io.github.cdimascio.dotenv.Dotenv;

@SpringBootApplication
@EnableAsync
public class AuthServiceApplication {

	public static void main(String[] args) {
		// Load .env và set system properties
		Dotenv dotenv = Dotenv.configure()
				.directory("../../")   // tùy vị trí .env
				.ignoreIfMissing()
				.load();

		dotenv.entries().forEach(entry -> System.setProperty(entry.getKey(), entry.getValue()));
		
		SpringApplication.run(AuthServiceApplication.class, args);
	}

}
