package club.tp0t.oj.Graphql.resolvers;

import com.coxautodev.graphql.tools.GraphQLQueryResolver;
import org.springframework.stereotype.Component;

@Component
public class TestResolver implements GraphQLQueryResolver {
    public String test() {
        return "hello world";
    }
}
