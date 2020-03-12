package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Graphql.types.*;
import club.tp0t.oj.Service.*;
import club.tp0t.oj.Util.CompetitionHelper;
import club.tp0t.oj.Util.RankHelper;
import com.coxautodev.graphql.tools.GraphQLMutationResolver;
import graphql.schema.DataFetchingEnvironment;
import graphql.servlet.context.DefaultGraphQLServletContext;
import org.springframework.stereotype.Component;

import javax.servlet.http.HttpSession;

@Component
public class AdminMutation implements GraphQLMutationResolver {
    private final BulletinService bulletinService;
    private final ChallengeService challengeService;
    private final UserService userService;
    private final RankHelper rankHelper;
    private final CompetitionHelper competitionHelper;

    public AdminMutation(BulletinService bulletinService, ChallengeService challengeService, UserService userService, RankHelper rankHelper, CompetitionHelper competitionHelper) {
        this.bulletinService = bulletinService;
        this.challengeService = challengeService;
        this.userService = userService;
        this.rankHelper = rankHelper;
        this.competitionHelper = competitionHelper;
    }

    public BulletinPubResult bulletinPub(BulletinPubInput bulletinPubInput, DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // login & admin check
        if (session.getAttribute("isLogin") == null ||
                !((boolean) session.getAttribute("isLogin")) ||
                !(boolean) session.getAttribute("isAdmin")) {
            return new BulletinPubResult("forbidden");
        }

        // input format check
        if (!bulletinPubInput.checkPass()) return new BulletinPubResult("not empty error");

        // unpack input data
        String title = bulletinPubInput.getTitle();
        String content = bulletinPubInput.getContent();
        boolean topping = bulletinPubInput.getTopping();

        // execute
        if (bulletinService.addBulletin(title, content, topping)) {
            return new BulletinPubResult("");
        } else {
            return new BulletinPubResult("Bulletin addition failed!");
        }
    }

    public ChallengeMutateResult challengeMutate(ChallengeMutateInput challengeMutate, DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // login & admin check
        if (session.getAttribute("isLogin") == null ||
                !((boolean) session.getAttribute("isLogin")) ||
                !((boolean) session.getAttribute("isAdmin") || (boolean) session.getAttribute("isTeam"))) {
            return new ChallengeMutateResult("forbidden");
        }

        // input format check
        if (!challengeMutate.checkPass()) return new ChallengeMutateResult("Challenge format not available");

        // unpack input data
        String id = challengeMutate.getChallengeId();

        // execute
        if (id.equals("")) {
            return new ChallengeMutateResult(challengeService.addChallenge(challengeMutate));
        } else {
            return new ChallengeMutateResult(challengeService.updateChallenge(challengeMutate));
        }
    }

    public ChallengeRemoveResult challengeRemove(String id) {
        return new ChallengeRemoveResult("Can't remove any challenge");
        // TODO:
//        if(!challengeService.removeById(id)) return new ChallengeRemoveResult("Remove failed");
//        return new ChallengeRemoveResult("");
    }

    public UserInfoUpdateResult userInfoUpdate(UserInfoUpdateInput input, DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // login & admin check
        if (session.getAttribute("isLogin") == null ||
                !((boolean) session.getAttribute("isLogin")) ||
                !(boolean) session.getAttribute("isAdmin")) {
            return new UserInfoUpdateResult("forbidden");
        }

        // input format check
        if (!input.checkPass()) return new UserInfoUpdateResult("");

        // cannot change one's own role
        long userId = (long) session.getAttribute("userId");
        if (userService.adminCheckByUserId(userId) &&
                Long.parseLong(input.getUserId()) == userId && !input.getRole().equals("admin")) {
            return new UserInfoUpdateResult("downgrade not permitted");
        }

        // execute
        userService.updateUserInfo(input.getUserId(),
                input.getName(),
                input.getRole(),
                input.getDepartment(),
                input.getGrade(),
                input.getProtectedTime(),
                input.getQq(),
                input.getMail(),
                input.getState());

        return new UserInfoUpdateResult("");
    }

    public boolean warmUp(DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // login & admin check
        if (session.getAttribute("isLogin") == null ||
                !((boolean) session.getAttribute("isLogin")) ||
                !(boolean) session.getAttribute("isAdmin")) {
            return false;
        }

        return rankHelper.warmUp();
    }

    public CompetitionMutationResult competition(CompetitionMutationInput input, DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // login & admin check
        if (session.getAttribute("isLogin") == null ||
                !((boolean) session.getAttribute("isLogin")) ||
                !(boolean) session.getAttribute("isAdmin")) {
            return new CompetitionMutationResult("forbidden");
        }


        CompetitionMutationResult competitionMutationResult = new CompetitionMutationResult("");
        competitionHelper.setCompetition(input.getCompetition());
        competitionHelper.setRegisterAllow(input.getRegisterAllow());
        competitionHelper.setBeginTime(input.getBeginTime());
        competitionHelper.setEndTime(input.getEndTime());
        return competitionMutationResult;
    }
}

