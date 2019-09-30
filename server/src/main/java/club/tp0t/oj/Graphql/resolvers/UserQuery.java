package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Entity.User;
import club.tp0t.oj.Graphql.types.UsersPayload;
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

    // get user info
    public UsersPayload users() {
        // TODO: session check
        // not login yet
        if(false) {
            return new UsersPayload("forbidden");
        }

        // normal user
        // can only get normal users
        if(false) {
            UsersPayload usersPayload = new UsersPayload("normal");
            List<User> users = userService.getNormalUsers();
            usersPayload.addNormalUserInfo(users);
            return usersPayload;
        }

        // admin
        // get admin users, disabled users
        UsersPayload usersPayload = new UsersPayload("admin");
        List<User> users = userService.getAllUsers();
        usersPayload.addAllUserInfo(users);
        return usersPayload;

    }
}
