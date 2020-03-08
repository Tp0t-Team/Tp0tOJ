package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.Submit;
import club.tp0t.oj.Graphql.types.*;
import club.tp0t.oj.Service.*;
import club.tp0t.oj.Util.ChallengeConfiguration;
import club.tp0t.oj.Util.OjConfig;
import com.coxautodev.graphql.tools.GraphQLQueryResolver;
import graphql.schema.DataFetchingEnvironment;
import graphql.servlet.context.DefaultGraphQLServletContext;
import org.springframework.stereotype.Component;

import javax.servlet.http.HttpSession;
import java.util.ArrayList;
import java.util.List;

@Component
public class AdminQuery implements GraphQLQueryResolver {
    private final ChallengeService challengeService;
    private final SubmitService submitService;
    private final UserService userService;
    private final OjConfig config;

    public AdminQuery(ChallengeService challengeService, SubmitService submitService, UserService userService, OjConfig config) {
        this.challengeService = challengeService;
        this.submitService = submitService;
        this.userService = userService;
        this.config = config;
    }

    public AllUserInfoResult allUserInfos(DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // login & admin check
        if (session.getAttribute("isLogin") == null ||
                !((boolean) session.getAttribute("isLogin")) ||
                !(boolean) session.getAttribute("isAdmin")) {
            return new AllUserInfoResult("forbidden");
        }

        // execute
        AllUserInfoResult res = new AllUserInfoResult("");
        res.addAllUserInfo(userService.getAllUser());

        return res;
    }

    public ChallengeConfigsResult challengeConfigs(DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // login & admin check
        if (session.getAttribute("isLogin") == null ||
                !((boolean) session.getAttribute("isLogin")) ||
                !((boolean) session.getAttribute("isAdmin") || (boolean) session.getAttribute("isTeam"))) {
            return new ChallengeConfigsResult("forbidden");
        }

        // execute
        List<Challenge> challenges = challengeService.getAllChallenges();

        // pack result
        List<ChallengeConfig> challengeConfigs = new ArrayList<>();
        for (Challenge challenge : challenges) {
            ChallengeConfiguration challengeconfiguration = ChallengeConfiguration.parseConfiguration(challenge.getConfiguration());
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
        ChallengeConfigsResult res = new ChallengeConfigsResult("");
        res.setChallengeConfigs(challengeConfigs);

        return res;
    }

    public SubmitHistoryResult submitHistory(String userId, DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        if (!config.isCompetition()) {
            // login & admin check
            if (session.getAttribute("isLogin") == null ||
                    !((boolean) session.getAttribute("isLogin")) ||
                    !(boolean) session.getAttribute("isAdmin")) {
                return new SubmitHistoryResult("forbidden");
            }
        } else {
            if (session.getAttribute("isLogin") == null) {
                return new SubmitHistoryResult("forbidden");
            }
        }

        // unpack input data
        long parsedUserId = Long.parseLong(userId);

        // execute
        List<Submit> submits = submitService.getCorrectSubmitsByUserId(parsedUserId);

        // pack result
        if (submits == null) return new SubmitHistoryResult("No such user.");
        SubmitHistoryResult submitHistoryResult = new SubmitHistoryResult("");
        submitHistoryResult.addSubmitInfos(submits);
        return submitHistoryResult;
    }
}
