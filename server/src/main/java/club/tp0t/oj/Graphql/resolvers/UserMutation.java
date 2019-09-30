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
    public RegisterPayload register(RegisterInput registerInput) {

        String name = registerInput.getName();
        String stuNumber = registerInput.getStuNumber();
        String password = registerInput.getPassword();
        String department = registerInput.getDepartment();
        String qq = registerInput.getQq();
        String mail = registerInput.getMail();

        // not empty
        if(name==null ||
                stuNumber==null ||
                password==null ||
                department==null ||
                qq==null ||
                mail==null) {
            return new RegisterPayload("not empty error");
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
            return new RegisterPayload("invalid mail");
        }

        // check duplicated user
        if(userService.checkNameExistence(name)) {
            return new RegisterPayload("name already in use");
        }
        if(userService.checkStuNumberExistence(stuNumber)) {
            return new RegisterPayload("stuNumber already in use");
        }
        if(userService.checkQqExistence(qq)) {
            return new RegisterPayload("qq already in use");
        }
        if(userService.checkMailExistence(mail)) {
            return new RegisterPayload("mail already in use");
        }

        // register user
        // if succeeded
        if(userService.register(name, stuNumber, password, department, qq, mail)) {
            return new RegisterPayload("success");
        }

        // if failed
        return new RegisterPayload("failed");
    }

    // user password reset
    // currently disabled
    public ResetPayload reset(ResetInput input) {

        // TODO: validate user info

        // TODO: reset password

        // if succeeded
        return new ResetPayload("success");

    }

    // user login
    public LoginPayload login(LoginInput input, DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // login check
        String name = input.getName();
        String password = input.getPassword();
        // not empty
        if(name==null) return new LoginPayload("failed");
        name = name.replaceAll("\\s", "");


        // user not exists
        if(!userService.checkNameExistence(name)) {
            return new LoginPayload("failed");
        }

        // login succeeded
        if(userService.login(name, password)) {
            session.setAttribute("isLogin", true);
            // admin
            if(userService.adminCheck(name)) {
                session.setAttribute("isAdmin", true);
            }
            else session.setAttribute("isAdmin", false);

            System.out.println("login succeeded");
            return new LoginPayload("success", userService.getIdByName(name));
        }
        // login failed
        else {
            System.out.println("login failed");
            return new LoginPayload("failed");
        }
    }

}
