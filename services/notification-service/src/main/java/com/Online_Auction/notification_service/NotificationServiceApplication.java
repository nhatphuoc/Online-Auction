package com.Online_Auction.notification_service;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

import org.springframework.amqp.rabbit.annotation.EnableRabbit;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
@EnableRabbit
public class NotificationServiceApplication {

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

		SpringApplication.run(NotificationServiceApplication.class, args);
	}

}
