package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Entity.User;
import club.tp0t.oj.Graphql.types.ChallengeInfosResult;
import club.tp0t.oj.Graphql.types.CompetitionResult;
import club.tp0t.oj.Graphql.types.RankResult;
import club.tp0t.oj.Graphql.types.UserInfoResult;
import club.tp0t.oj.Service.ChallengeService;
import club.tp0t.oj.Service.UserService;
import com.coxautodev.graphql.tools.GraphQLQueryResolver;
import graphql.schema.DataFetchingEnvironment;
import graphql.servlet.context.DefaultGraphQLServletContext;
import org.springframework.stereotype.Component;
import club.tp0t.oj.Service.ComputationService;
import javax.servlet.http.HttpSession;

@Component
public class UserQuery implements GraphQLQueryResolver {
    private final ChallengeService challengeService;
    private final UserService userService;
    private final ComputationService computationService;
    public UserQuery(ChallengeService challengeService, UserService userService,ComputationService computationService) {
        this.challengeService = challengeService;
        this.userService = userService;
        this.computationService = computationService;
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
            userInfoResult.addOwnUserInfo(user, userService.getRankByUserId(user.getUserId()));

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
                userInfoResult.addOthersUserInfo(user, userService.getRankByUserId(user.getUserId()));
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

        // execute
        ChallengeInfosResult challengeInfosResult = new ChallengeInfosResult("");
        challengeInfosResult.setChallengeInfos(challengeService.getChallengeInfoForUser((long) session.getAttribute("userId")));

        return challengeInfosResult;
    }

    public CompetitionResult competition(DataFetchingEnvironment environment) {
        //TODO:

        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        CompetitionResult competitionResult = new CompetitionResult("");
        competitionResult.setCompetitionMode(computationService.getCompetitionMode());
        competitionResult.setRegisterAllow(computationService.getRegisterAllow());
        competitionResult.setBeginTime(computationService.getBeginTime());
        competitionResult.setBeginTime(computationService.getEndTime());

        return competitionResult;
    }

}
