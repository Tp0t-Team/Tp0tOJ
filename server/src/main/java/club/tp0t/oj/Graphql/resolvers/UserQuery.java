package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Entity.User;
import club.tp0t.oj.Graphql.types.ChallengeInfosResult;
import club.tp0t.oj.Graphql.types.CompetitionResult;
import club.tp0t.oj.Graphql.types.RankResult;
import club.tp0t.oj.Graphql.types.UserInfoResult;
import club.tp0t.oj.Service.ChallengeService;
import club.tp0t.oj.Service.UserService;
import club.tp0t.oj.Util.OjConfig;
import club.tp0t.oj.Util.RankHelper;
import com.coxautodev.graphql.tools.GraphQLQueryResolver;
import graphql.schema.DataFetchingEnvironment;
import graphql.servlet.context.DefaultGraphQLServletContext;
import org.springframework.stereotype.Component;
import club.tp0t.oj.Util.CompetitionHelper;

import javax.servlet.http.HttpSession;
import java.util.Date;

@Component
public class UserQuery implements GraphQLQueryResolver {
    private final ChallengeService challengeService;
    private final UserService userService;
    private final CompetitionHelper competitionHelper;
    private final RankHelper rankHelper;
    private final OjConfig config;

    public UserQuery(ChallengeService challengeService, UserService userService, CompetitionHelper competitionHelper, RankHelper rankHelper, OjConfig config) {
        this.challengeService = challengeService;
        this.userService = userService;
        this.competitionHelper = competitionHelper;
        this.rankHelper = rankHelper;
        this.config = config;
    }

    // test
    public String test(DataFetchingEnvironment environment) {
        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();
        session.setAttribute("test", true);
        return "hello world";
    }

    // get user rank
    public RankResult rank(DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // not login yet
        //if(session.getAttribute("isLogin") == null  || !(boolean)session.getAttribute("isLogin")) {
        //    session.setAttribute("isLogin", false);
        //    return new RankResult("forbidden");
        //}

        //execute
        RankResult rankResult = new RankResult("");
        rankResult.addRankResultDescs(userService.getUsersRank());

        return rankResult;
    }

    // get user profile
    // TODO: is necessary moving into service?
    public UserInfoResult userInfo(String userId, DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // if not login
        if (session.getAttribute("isLogin") == null || !(boolean) session.getAttribute("isLogin")) {
            return new UserInfoResult("forbidden");
        }

        // unpack input data
        long parsedUserId = Long.parseLong(userId);
        long currentUserId = (Long) session.getAttribute("userId");

        // if requested by himself or by admin
        if (currentUserId == parsedUserId ||
                userService.adminCheckByUserId(currentUserId)) {
            User user = userService.getUserById(parsedUserId);
            UserInfoResult userInfoResult = new UserInfoResult("");
            userInfoResult.addOwnUserInfo(user, rankHelper.getUserInfo(user.getUserId()));

            return userInfoResult;
        }
        // if requested by other users
        else {
            // if user not exists
            if (!userService.checkUserIdExistence(parsedUserId)) {
                return new UserInfoResult("not found");
            }
            // if user exists
            else {
                User user = userService.getUserById(parsedUserId);
                UserInfoResult userInfoResult = new UserInfoResult("");
                userInfoResult.addOthersUserInfo(user, rankHelper.getUserInfo(user.getUserId()));
                return userInfoResult;
            }
        }
    }

    // get challenges
    public ChallengeInfosResult challengeInfos(DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // if not login
        if (session.getAttribute("isLogin") == null || !(boolean) session.getAttribute("isLogin")) {
            return new ChallengeInfosResult("forbidden");
        }

        if (competitionHelper.getCompetition()) {
            Date now = new Date();
            if (now.compareTo(competitionHelper.getBeginTime()) < 0) {
                return new ChallengeInfosResult("");
            }
        }

        // execute
        ChallengeInfosResult challengeInfosResult = new ChallengeInfosResult("");
        challengeInfosResult.setChallengeInfos(challengeService.getChallengeInfoForUser((long) session.getAttribute("userId")));

        return challengeInfosResult;
    }

    public CompetitionResult competition(DataFetchingEnvironment environment) {

        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        CompetitionResult competitionResult = new CompetitionResult("");
        competitionResult.setCompetition(competitionHelper.getCompetition());
        competitionResult.setRegisterAllow(competitionHelper.getRegisterAllow());
        competitionResult.setBeginTime(config.getBeginTime());
        competitionResult.setEndTime(config.getEndTime());
        return competitionResult;
    }

}
