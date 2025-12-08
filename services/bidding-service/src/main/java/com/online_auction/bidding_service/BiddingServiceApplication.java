package com.online_auction.bidding_service;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import io.github.cdimascio.dotenv.Dotenv;

@SpringBootApplication
public class BiddingServiceApplication {

	public static void main(String[] args) {
		// Load .env và set system properties
		Dotenv dotenv = Dotenv.configure()
				.directory("../../")   // tùy vị trí .env
				.ignoreIfMissing()
				.load();

		dotenv.entries().forEach(entry -> System.setProperty(entry.getKey(), entry.getValue()));

		SpringApplication.run(BiddingServiceApplication.class, args);
	}

}
