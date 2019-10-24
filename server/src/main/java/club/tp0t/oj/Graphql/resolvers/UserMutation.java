package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.Replica;
import club.tp0t.oj.Entity.ReplicaAlloc;
import club.tp0t.oj.Entity.User;
import club.tp0t.oj.Graphql.types.*;
import club.tp0t.oj.Service.*;
import club.tp0t.oj.Util.ChallengeConfiguration;
import club.tp0t.oj.Util.CheckHelper;
import com.coxautodev.graphql.tools.GraphQLMutationResolver;
import graphql.schema.DataFetchingEnvironment;
import graphql.servlet.context.DefaultGraphQLServletContext;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Transactional;

import javax.servlet.http.HttpSession;
import java.util.ArrayList;
import java.util.List;
import java.util.regex.Pattern;

@Component
public class UserMutation implements GraphQLMutationResolver {
    @Autowired
    private BulletinService bulletinService;
    @Autowired
    private ChallengeService challengeService;
    @Autowired
    private FlagService flagService;
    @Autowired
    private SubmitService submitService;
    @Autowired
    private UserService userService;

    // user register
    public RegisterResult register(RegisterInput registerInput, DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // already login
        if ((session.getAttribute("isLogin") != null && (boolean) session.getAttribute("isLogin"))) {
            return new RegisterResult("already login cannot register");
        }

        // unpack input data
        String name = registerInput.getName();
        String stuNumber = registerInput.getStuNumber();
        String password = registerInput.getPassword();
        String department = registerInput.getDepartment();
        String qq = registerInput.getQq();
        String mail = registerInput.getMail();
        String grade = registerInput.getGrade();

        // input format check
        if (!registerInput.checkPass()) return new RegisterResult("invalid information");
        if (!CheckHelper.MAIL_PATTERN.matcher(registerInput.getMail()).matches()) {
            return new RegisterResult("invalid mail");
        }
        // TODO: validate you are a student

        // execute
        return new RegisterResult(userService.register(name, stuNumber, password, department, qq, mail, grade));
    }

    // user password reset
    // currently disabled
    /*
    public ResetResult reset(ResetInput input) {

        // validate user info

        // reset password

        // if succeeded
        return new ResetResult("success");

    }
    */

    // user login
    public LoginResult login(LoginInput input, DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // already login
        //if((session.getAttribute("isLogin") != null && (boolean)session.getAttribute("isLogin"))) {
        //    return new LoginResult("already login");
        //}

        // input format check
        if (!input.checkPass()) return new LoginResult("not empty error");

        // unpack input data
        String stuNumber = input.getStuNumber();
        String password = input.getPassword();

        // execute
        User user = userService.login(stuNumber, password);
        if (user != null) {
            session.setAttribute("isLogin", true);
            session.setAttribute("userId", user.getUserId());
            // admin
            if (user.getRole().equals("admin")) {
                session.setAttribute("isAdmin", true);
            } else session.setAttribute("isAdmin", false);
            // team
            if (user.getRole().equals("team")) {
                session.setAttribute("isTeam", true);
            } else session.setAttribute("isTeam", false);

            return new LoginResult("", Long.toString(user.getUserId()), user.getRole());
        } else {
            return new LoginResult("failed");
        }
    }

    // user logout
    public LogoutResult logout(DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        if ((session.getAttribute("isLogin") == null || !(boolean) session.getAttribute("isLogin"))) {
            return new LogoutResult("not login yet");
        }

        session.setAttribute("isLogin", false);
        return new LogoutResult("");
    }

    // submit flag
    public SubmitResult submit(SubmitInput input, DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // not login yet
        if (session.getAttribute("isLogin") == null || !(boolean) session.getAttribute("isLogin")) {
            return new SubmitResult("forbidden");
        }

        // input format check
        if (!input.checkPass()) return new SubmitResult("not empty error");

        // unpack input data
        long challengeId = Long.parseLong(input.getChallengeId());
        long userId = (long) session.getAttribute("userId");
        String flag = input.getFlag();
        boolean isMember = !(boolean) session.getAttribute("isAdmin") && !(boolean) session.getAttribute("isTeam");

        // execute
        return new SubmitResult(submitService.submit(userId, challengeId, flag, isMember));
    }
}
