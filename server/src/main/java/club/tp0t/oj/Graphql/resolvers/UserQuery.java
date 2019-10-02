package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.User;
import club.tp0t.oj.Graphql.types.ChallengesResult;
import club.tp0t.oj.Graphql.types.Profile;
import club.tp0t.oj.Graphql.types.RankResult;
import club.tp0t.oj.Service.*;
import graphql.schema.DataFetchingEnvironment;
import graphql.servlet.context.DefaultGraphQLServletContext;
import org.springframework.beans.factory.annotation.Autowired;

import javax.servlet.http.HttpSession;
import java.util.List;

public class UserQuery extends Query {
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
        if(session.getAttribute("isLogin") == null  || !(boolean)session.getAttribute("isLogin")) {
            session.setAttribute("isLogin", false);
            return new RankResult("forbidden");
        }

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

        rankResult.addUserInfos(users);
        return rankResult;
    }

    // get user profile
    public Profile userInfo(String userId, DataFetchingEnvironment environment) {
        // TODO: not implemented
        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // if not login
        if(session.getAttribute("isLogin")==null || !(boolean)session.getAttribute("isLogin")) {
            return new Profile("forbidden");
        }

        // whether requested by user himself
        long currentUserId = (Long) session.getAttribute("userId");
        System.out.println("currentUserId: " + currentUserId);
        // by himself
        if(currentUserId == Long.parseLong(userId)) {
            User user = userService.getUserById(Long.parseLong(userId));
            System.out.println(user);
            Profile profile = new Profile("requested by self");
            profile.addOwnUserInfo(user, userService.getRankByUserId(user.getUserId()));
            return profile;

        }
        // if requested by other users
        else {
            User user = userService.getUserById(Long.parseLong(userId));
            // user not exists
            if(user == null) {
                return new Profile("not found");
            }
            // exists
            else {
                Profile profile = new Profile("requested by others");
                profile.addOthersUserInfo(user, userService.getRankByUserId(user.getUserId()));
                return profile;
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
        System.out.println("have challenges: " + challenges.size());
        System.out.println("userId: " + (long)session.getAttribute("userId"));
        challengesResult.addChallengeInfos(challenges, (long)session.getAttribute("userId"), submitService);

        return challengesResult;
    }

}
