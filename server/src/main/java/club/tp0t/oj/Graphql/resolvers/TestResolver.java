package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Dao.*;
import com.coxautodev.graphql.tools.GraphQLQueryResolver;
import graphql.schema.DataFetchingEnvironment;
import graphql.servlet.context.DefaultGraphQLServletContext;
import org.springframework.stereotype.Component;

import javax.servlet.http.HttpSession;

@Component
public class TestResolver implements GraphQLQueryResolver {
    private final BulletinRepository bulletinRepository;
    private final ChallengeRepository challengeRepository;
    private final FlagRepository flagRepository;
    private final ReplicaRepository replicaRepository;
    private final ReplicaAllocRepository replicaAllocRepository;
    private final SubmitRepository submitRepository;
    private final UserRepository userRepository;

    public TestResolver(BulletinRepository bulletinRepository,
                        ChallengeRepository challengeRepository,
                        FlagRepository flagRepository,
                        ReplicaRepository replicaRepository,
                        ReplicaAllocRepository replicaAllocRepository,
                        SubmitRepository submitRepository,
                        UserRepository userRepository) {
        this.bulletinRepository = bulletinRepository;
        this.challengeRepository = challengeRepository;
        this.flagRepository = flagRepository;
        this.replicaAllocRepository = replicaAllocRepository;
        this.replicaRepository = replicaRepository;
        this.submitRepository = submitRepository;
        this.userRepository = userRepository;
    }

    // test
    public String test(DataFetchingEnvironment environment) {
        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();
        session.setAttribute("test", true);
        return "hello world";
    }
}
