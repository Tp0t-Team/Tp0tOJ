package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Entity.User;
import club.tp0t.oj.Graphql.types.*;
import club.tp0t.oj.Service.SubmitService;
import club.tp0t.oj.Service.UserService;
import club.tp0t.oj.Util.CheckHelper;
import club.tp0t.oj.Util.OjConfig;
import com.coxautodev.graphql.tools.GraphQLMutationResolver;
import graphql.schema.DataFetchingEnvironment;
import graphql.servlet.context.DefaultGraphQLServletContext;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.mail.SimpleMailMessage;
import org.springframework.mail.javamail.JavaMailSender;
import org.springframework.mail.javamail.MimeMailMessage;
import org.springframework.mail.javamail.MimeMessageHelper;
import org.springframework.stereotype.Component;
import org.springframework.util.DigestUtils;

import javax.mail.internet.MimeMessage;
import javax.mail.internet.MimeUtility;
import javax.servlet.http.HttpSession;
import java.util.Date;

@Component
public class UserMutation implements GraphQLMutationResolver {
    private final SubmitService submitService;
    private final UserService userService;
    private final JavaMailSender mailSender;
    private final OjConfig config;

    @Value("${spring.mail.username}")
    private String from;

    @Autowired
    public UserMutation(SubmitService submitService, UserService userService, JavaMailSender mailSender, OjConfig config) {
        this.submitService = submitService;
        this.userService = userService;
        this.mailSender = mailSender;
        this.config = config;
    }

    // user register
    public RegisterResult register(RegisterInput registerInput, DataFetchingEnvironment environment) {

        if(!config.isAllowRegister()) {
            return new RegisterResult("disabled");
        }

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
        String password = DigestUtils.md5DigestAsHex((config.getSalt() + registerInput.getPassword()).getBytes()); // registerInput.getPassword();
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
        String password = DigestUtils.md5DigestAsHex((config.getSalt() + input.getPassword()).getBytes()); // input.getPassword();

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

    public ForgetResult forget(String input) {

        if(!config.isAllowRegister()) {
            return new ForgetResult("disabled");
        }

        // input format check
        if (input.equals("")) return new ForgetResult("not empty error");

        //execute
        UserService.Forget forget = userService.forget(input);
        if (!forget.getMessage().equals("")) {
            return new ForgetResult(forget.getMessage());
        } else {
            // TODO: send mail
//            System.out.println(forget.getToken()); // for debug
            try {
                MimeMessageHelper mailMessage = new MimeMessageHelper(mailSender.createMimeMessage());
                String toName = MimeUtility.encodeText(forget.getName());
                String fromName = MimeUtility.encodeText(config.getName());
                mailMessage.setTo(String.format("%s <%s>", toName, forget.getMail()));
                mailMessage.setFrom(String.format("%s <%s>", fromName, from));
                mailMessage.setSubject("password reset");
                mailMessage.setText(String.format("Please use the follow link to reset your password.\n%s/reset?token=%s", config.getHost(), forget.getToken()));
                mailSender.send(mailMessage.getMimeMessage());
                return new ForgetResult("");
            } catch (Exception e) {
                return new ForgetResult("Send mail failed, please wait.");
            }
        }
    }

    public ResetResult reset(ResetInput input) {

        // input format check
        if (!input.checkPass()) return new ResetResult("not empty error");

        // unpack input data
        String token = input.getToken();
        String password = DigestUtils.md5DigestAsHex((config.getSalt() + input.getPassword()).getBytes()); // input.getPassword();

        // execute
        return new ResetResult(userService.reset(token, password));
    }

    // submit flag
    public SubmitResult submit(SubmitInput input, DataFetchingEnvironment environment) {

        Date now = new Date();
        if(now.compareTo(config.getBeginTime()) < 0) {
            return new SubmitResult("disabled");
        }else if(now.compareTo(config.getEndTime()) > 0) {
            return new SubmitResult("competition finished");
        }

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
