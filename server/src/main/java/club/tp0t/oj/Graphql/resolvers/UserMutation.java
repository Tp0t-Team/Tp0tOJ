package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Graphql.types.*;
import club.tp0t.oj.Service.*;
import com.coxautodev.graphql.tools.GraphQLMutationResolver;
import graphql.schema.DataFetchingEnvironment;
import graphql.servlet.context.DefaultGraphQLServletContext;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import javax.servlet.http.HttpSession;
import java.util.regex.Pattern;

@Component
public class UserMutation implements GraphQLMutationResolver {
    @Autowired
    private  BulletinService bulletinService;
    @Autowired
    private  ChallengeService challengeService;
    @Autowired
    private  FlagService flagService;
    @Autowired
    private  ReplicaService replicaService;
    @Autowired
    private  ReplicaAllocService replicaAllocService;
    @Autowired
    private  SubmitService submitService;
    @Autowired
    private  UserService userService;

    // user register
    public RegisterResult register(RegisterInput registerInput) {

        String name = registerInput.getName();
        String stuNumber = registerInput.getStuNumber();
        String password = registerInput.getPassword();
        String department = registerInput.getDepartment();
        String qq = registerInput.getQq();
        String mail = registerInput.getMail();
        String grade = registerInput.getGrade();

        // not empty
        if(name==null ||
                stuNumber==null ||
                password==null ||
                department==null ||
                qq==null ||
                mail==null ||
                grade==null) {
            return new RegisterResult("not empty error");
        }

        name = name.replaceAll("\\s", "");
        stuNumber = stuNumber.replaceAll("\\s", "");
        qq = qq.replaceAll("\\s", "");
        mail = mail.replaceAll("\\s", "");


        // TODO: validate user info
        // validate email
        String EMAIL_PATTERN = "^[_A-Za-z0-9-+]+(.[_A-Za-z0-9-]+)*@" +
                "[A-Za-z0-9-]+(.[A-Za-z0-9]+)*(.[A-Za-z]{2,})$";
        Pattern pattern = Pattern.compile(EMAIL_PATTERN);
        if(!pattern.matcher(registerInput.getMail()).matches()) {
            return new RegisterResult("invalid mail");
        }

        // check duplicated user
        if(userService.checkStuNumberExistence(stuNumber)) {
            return new RegisterResult("stuNumber already in use");
        }
        if(userService.checkQqExistence(qq)) {
            return new RegisterResult("qq already in use");
        }
        if(userService.checkMailExistence(mail)) {
            return new RegisterResult("mail already in use");
        }

        // register user
        // if succeeded
        if(userService.register(name, stuNumber, password, department, qq, mail, grade)) {
            return new RegisterResult("success");
        }

        // if failed
        return new RegisterResult("failed");
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

        // login check
        String stuNumber = input.getStuNumber();
        String password = input.getPassword();
        // not empty
        if(stuNumber==null) return new LoginResult("not empty error");
        stuNumber = stuNumber.replaceAll("\\s", "");
        if(stuNumber.equals("")) return new LoginResult("not empty error");

        // user not exists
        if(!userService.checkStuNumberExistence(stuNumber)) {
            return new LoginResult("failed");
        }

        // user password check succeeded
        if(userService.login(stuNumber, password)) {
            session.setAttribute("isLogin", true);
            session.setAttribute("userId", userService.getIdByStuNumber(stuNumber));
            // admin
            if(userService.adminCheckByStuNumber(stuNumber)) {
                session.setAttribute("isAdmin", true);
            }
            else session.setAttribute("isAdmin", false);
            // team
            if(userService.teamCheckByStuNumber(stuNumber)) {
                session.setAttribute("isTeam", true);
            }
            else session.setAttribute("isTeam", false);

            System.out.println("login succeeded");
            return new LoginResult("success", Long.toString(userService.getIdByStuNumber(stuNumber)),
                    userService.getRoleByStuNumber(stuNumber));
        }
        // user password check failed
        else {
            System.out.println("login failed");
            return new LoginResult("failed");
        }
    }

    // user logout
    public LogoutResult logout(DataFetchingEnvironment environment) {
        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        session.setAttribute("isLogin", false);
        return new LogoutResult("success");
    }

    // submit flag
    /*
    public SubmitResult submit(SubmitInput input, DataFetchingEnvironment environment) {
        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // not login yet
        if(session.getAttribute("isLogin") == null || !(boolean)session.getAttribute("isLogin")) {
            return new SubmitResult("forbidden");
        }

        // not empty
        if(input.getChallengeId() == null || input.getFlag() == null) {
            return new SubmitResult("not empty error");
        }

        // check flag
        long challengeId = Long.parseLong(input.getChallengeId());
        long userId = (long) session.getAttribute("userId");
        String flag = flagService.getFlagByUserIdAndChallengeId(userId, challengeId);
        String submitFlag = input.getFlag();
        // correct flag
        boolean correct = false;
        int mark = 0;
        if(submitFlag.equals(flag)) {
            // TODO: duplicate submit
            if(submitService.checkDuplicateSubmit(userService.getUserById(userId), challengeId)) {
                return new SubmitResult("duplicate submit");
            }

            correct = true;
            // add to user score
            // TODO: get points of challenge
            long points = 100;
            userService.addScore(points);

            // whether first three solvers

        }
        else {

        }

        // save into submit table
        submitService.submit(userService.getUserById(userId), submitFlag, correct, mark);
    }
    */

}
