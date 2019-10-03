package club.tp0t.oj;

import club.tp0t.oj.Graphql.resolvers.AdminMutation;
import club.tp0t.oj.Graphql.resolvers.AdminQuery;
import club.tp0t.oj.Graphql.resolvers.UserMutation;
import club.tp0t.oj.Graphql.resolvers.UserQuery;
import club.tp0t.oj.Graphql.resolvers.BulletinQuery;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.autoconfigure.http.HttpMessageConverters;
import org.springframework.context.annotation.Bean;
import org.springframework.http.MediaType;
import org.springframework.http.converter.HttpMessageConverter;

import java.util.ArrayList;
import java.util.List;

@SpringBootApplication
public class OjApplication {

	public static void main(String[] args) {
		SpringApplication.run(OjApplication.class, args);
	}

}
