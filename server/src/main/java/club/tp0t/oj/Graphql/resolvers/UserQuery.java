package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.Submit;
import club.tp0t.oj.Entity.User;
import club.tp0t.oj.Graphql.types.ChallengesResult;
import club.tp0t.oj.Graphql.types.SubmitHistoryResult;
import club.tp0t.oj.Graphql.types.UserInfoResult;
import club.tp0t.oj.Graphql.types.RankResult;
import club.tp0t.oj.Service.*;
import com.coxautodev.graphql.tools.GraphQLQueryResolver;
import graphql.schema.DataFetchingEnvironment;
import graphql.servlet.context.DefaultGraphQLServletContext;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import javax.servlet.http.HttpSession;
import java.util.List;

@Component
public class UserQuery implements GraphQLQueryResolver {
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
    private  SubmitService submitService;
    @Autowired
    private  UserService userService;

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

        /*
        // normal user
        // can only get normal users
        if(false) {
            Profile profile = new Profile("normal");
            List<User> users = userService.getNormalUsers();
            profile.addNormalUserInfo(users);
            return profile;
        }

        // admin
        // get admin users, disabled users
        Profile profile = new Profile("admin");
        List<User> users = userService.getAllUsers();
        profile.addAllUserInfo(users);
        return profile;
        */

        RankResult rankResult = new RankResult("");
        List<User> users = userService.getUsersRank();

        // no users
        if(users == null) return rankResult;

        rankResult.addRankResultDescs(users);
        return rankResult;
    }

    // get user profile
    public UserInfoResult userInfo(String userId, DataFetchingEnvironment environment) {
        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // if not login
        if(session.getAttribute("isLogin")==null || !(boolean)session.getAttribute("isLogin")) {
            return new UserInfoResult("forbidden");
        }

        // whether requested by user himself
        long currentUserId = (Long) session.getAttribute("userId");
        // by himself or by admin
        if(currentUserId == Long.parseLong(userId) ||
                userService.adminCheckByUserId(currentUserId)) {
            User user = userService.getUserById(Long.parseLong(userId));
            UserInfoResult userInfoResult = new UserInfoResult("");
            userInfoResult.addOwnUserInfo(user, userService.getRankByUserId(user.getUserId()));
            return userInfoResult;

        }
        // if requested by other users
        else {
            // user not exists
            if(!userService.checkUserIdExistence(Long.parseLong(userId))) {
                return new UserInfoResult("not found");
            }
            // exists
            else {
                User user = userService.getUserById(Long.parseLong(userId));
                UserInfoResult userInfoResult = new UserInfoResult("");
                userInfoResult.addOthersUserInfo(user, userService.getRankByUserId(user.getUserId()));
                return userInfoResult;
            }
        }

    }

    // get challenges
    public ChallengesResult challenges(DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // if not login
        if(session.getAttribute("isLogin")==null || !(boolean)session.getAttribute("isLogin")) {
            return new ChallengesResult("forbidden");
        }

        List<Challenge> challenges;

        // admin or team
        long userId = (long) session.getAttribute("userId");
        if(userService.adminCheckByUserId(userId) || userService.teamCheckByUserId(userId)) {
            challenges = challengeService.getAllChallenges();
        }
        // member
        else {
            challenges = challengeService.getEnabledChallenges();
        }

        // no challenge
        if(challenges == null) return new ChallengesResult("no challenge available");

        ChallengesResult challengesResult = new ChallengesResult("");
        challengesResult.addChallengeInfos(challenges, (long)session.getAttribute("userId"), submitService);

        return challengesResult;
    }

    // admin get user's submit history
    public SubmitHistoryResult submitHistory(String userId, DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // login & admin check
        if(session.getAttribute("isLogin") == null ||
                !((boolean) session.getAttribute("isLogin")) ||
                !(boolean) session.getAttribute("isAdmin")) {
            return new SubmitHistoryResult("forbidden");
        }

        User user = userService.getUserById(Long.parseLong(userId));
        List<Submit> submits = submitService.getCorrectSubmitsByUser(user);

        SubmitHistoryResult submitHistoryResult = new SubmitHistoryResult("");

        submitHistoryResult.addSubmitInfos(submits);

        return submitHistoryResult;
    }

}
