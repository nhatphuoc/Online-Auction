package com.Online_Auction.notification_service.config;

import org.springframework.amqp.core.*;
import org.springframework.amqp.rabbit.connection.ConnectionFactory;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.amqp.support.converter.Jackson2JsonMessageConverter;
import org.springframework.amqp.support.converter.MessageConverter;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class RabbitMQConfig {

    public static final String EXCHANGE = "auction.exchange";
    public static final String BID_SUCCESS_QUEUE = "bid.placed.success.queue";

    public static final String ROUTING_KEY_BID_SUCCESS = "bid.success";

    @Bean
    public DirectExchange exchange() {
        return new DirectExchange(EXCHANGE);
    }

    @Bean
    public Queue bidSuccessQueue() {
        return new Queue(BID_SUCCESS_QUEUE);
    }

    @Bean
    public Binding bindingBidSuccess(Queue bidSuccessQueue, DirectExchange exchange) {
        return BindingBuilder.bind(bidSuccessQueue).to(exchange).with(ROUTING_KEY_BID_SUCCESS);
    }

    @Bean public MessageConverter jsonMessageConverter() { 
        return new Jackson2JsonMessageConverter(); 
    }


}