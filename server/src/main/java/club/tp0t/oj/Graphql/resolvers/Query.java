package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Dao.*;
import com.coxautodev.graphql.tools.GraphQLQueryResolver;
import graphql.schema.DataFetchingEnvironment;
import graphql.servlet.context.DefaultGraphQLServletContext;
import org.springframework.stereotype.Component;

import javax.servlet.http.HttpSession;

@Component
public abstract class Query implements GraphQLQueryResolver {

}