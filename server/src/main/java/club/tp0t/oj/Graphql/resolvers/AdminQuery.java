package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.User;
import club.tp0t.oj.Graphql.types.AllUserInfoResult;
import club.tp0t.oj.Graphql.types.ChallengeConfig;
import club.tp0t.oj.Graphql.types.ChallengeConfigsResult;
import club.tp0t.oj.Graphql.types.UserInfoResult;
import club.tp0t.oj.Service.*;
import club.tp0t.oj.Util.ChallengeConfiguration;
import com.coxautodev.graphql.tools.GraphQLQueryResolver;
import graphql.schema.DataFetchingEnvironment;
import graphql.servlet.context.DefaultGraphQLServletContext;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import javax.servlet.http.HttpSession;
import java.util.ArrayList;
import java.util.List;

@Component
public class AdminQuery implements GraphQLQueryResolver {
    @Autowired
    private BulletinService bulletinService;
    @Autowired
    private ChallengeService challengeService;
    @Autowired
    private FlagService flagService;
    @Autowired
    private ReplicaService replicaService;
    @Autowired
    private ReplicaAllocService replicaAllocService;
    @Autowired
    private SubmitService submitService;
    @Autowired
    private UserService userService;

    public AllUserInfoResult allUserInfos(DataFetchingEnvironment environment) {
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // login & admin check
        if (session.getAttribute("isLogin") == null ||
                !((boolean) session.getAttribute("isLogin")) ||
                !(boolean) session.getAttribute("isAdmin")) {
            return new AllUserInfoResult("forbidden");
        }

        AllUserInfoResult res = new AllUserInfoResult("");

        List<User> users = userService.getAllUser();
        if(users == null) users = new ArrayList<>();
        res.addAllUserInfo(users);

        return res;
    }

    public ChallengeConfigsResult challengeConfigs(DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // login & admin check
        if (session.getAttribute("isLogin") == null ||
                !((boolean) session.getAttribute("isLogin")) ||
                !(boolean) session.getAttribute("isAdmin")) {
            return new ChallengeConfigsResult("forbidden");
        }

        List<ChallengeConfig> challengeConfigs = new ArrayList<>();
        ChallengeConfigsResult res = new ChallengeConfigsResult("");
        List<Challenge> challenges = challengeService.getAllChallenges();

        for (Challenge challenge : challenges) {
            ChallengeConfiguration challengeconfiguration = challengeService.getConfiguration(challenge);
            ChallengeConfig challengeConfig = new ChallengeConfig();

            challengeConfig.setChallengeId(Long.toString(challenge.getChallengeId()));
            challengeConfig.setState(challenge.getState());
            challengeConfig.setName(challengeconfiguration.getName());
            challengeConfig.setType(challengeconfiguration.getType());
            challengeConfig.setDescription(challengeconfiguration.getDescription());
            challengeConfig.setExternal_link(challengeconfiguration.getExternalLink());
            challengeConfig.setHint(challengeconfiguration.getHint());
            challengeConfig.setFlag(challengeconfiguration.getFlagEx());
            challengeConfig.setScore(challengeconfiguration.getScoreEx());

            challengeConfigs.add(challengeConfig);
        }
        res.setChallengeConfigs(challengeConfigs);

        return res;
    }
}
