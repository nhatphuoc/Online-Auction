package com.online_auction.bidding_service;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class BiddingServiceApplication {

	public static void main(String[] args) throws IOException {
		Path envPath = Paths.get("../../shared/.env");

		if (Files.exists(envPath)) {
			Files.lines(envPath)
					.map(String::trim)
					.filter(line -> !line.isEmpty()) // bỏ dòng trống
					.filter(line -> !line.startsWith("#")) // bỏ comment
					.filter(line -> line.matches("^[A-Za-z_][A-Za-z0-9_]*=.*")) // KEY=VALUE
					.forEach(line -> {
						String[] parts = line.split("=", 2);
						System.setProperty(parts[0], parts[1]);
					});
		}

		SpringApplication.run(BiddingServiceApplication.class, args);
	}

}
