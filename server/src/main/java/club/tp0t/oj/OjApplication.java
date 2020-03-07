package club.tp0t.oj;

import club.tp0t.oj.Graphql.resolvers.AdminMutation;
import club.tp0t.oj.Graphql.resolvers.AdminQuery;
import club.tp0t.oj.Graphql.resolvers.UserMutation;
import club.tp0t.oj.Graphql.resolvers.UserQuery;
import club.tp0t.oj.Graphql.resolvers.BulletinQuery;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.autoconfigure.http.HttpMessageConverters;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.boot.web.servlet.MultipartConfigFactory;
import org.springframework.context.annotation.Bean;
import org.springframework.data.jpa.repository.config.EnableJpaAuditing;
import org.springframework.http.MediaType;
import org.springframework.http.converter.HttpMessageConverter;
import org.springframework.util.unit.DataSize;

import javax.servlet.MultipartConfigElement;
import java.util.ArrayList;
import java.util.List;

@SpringBootApplication
@EnableJpaAuditing
@EnableConfigurationProperties
public class OjApplication {

    public static void main(String[] args) {
        SpringApplication.run(OjApplication.class, args);
    }

	/*
	@Bean
	public UserQuery userQuery() {
		return new UserQuery();
	}

	@Bean
	public AdminQuery adminQuery() {
		return new AdminQuery();
	}

	@Bean
	public BulletinQuery bulletinQuery() {
		return new BulletinQuery();
	}

	 */

    @Bean
    public MultipartConfigElement multipartConfigElement() {
        MultipartConfigFactory factory = new MultipartConfigFactory();
        factory.setMaxFileSize(DataSize.ofMegabytes(20));
        return factory.createMultipartConfig();
    }
}
