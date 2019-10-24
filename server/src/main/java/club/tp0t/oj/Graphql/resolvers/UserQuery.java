package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.User;
import club.tp0t.oj.Graphql.types.ChallengeInfosResult;
import club.tp0t.oj.Graphql.types.UserInfoResult;
import club.tp0t.oj.Graphql.types.RankResult;
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
public class UserQuery implements GraphQLQueryResolver {
    private final ChallengeService challengeService;
    private final UserService userService;

    public UserQuery(ChallengeService challengeService, UserService userService) {
        this.challengeService = challengeService;
        this.userService = userService;
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
        Long parsedUserId = Long.parseLong(userId);
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
            if (!userService.checkUserIdExistence(Long.parseLong(userId))) {
                return new UserInfoResult("not found");
            }
            // if user exists
            else {
                User user = userService.getUserById(Long.parseLong(userId));
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

}
