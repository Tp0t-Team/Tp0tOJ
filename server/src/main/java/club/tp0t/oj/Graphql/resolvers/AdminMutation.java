package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Graphql.types.*;
import club.tp0t.oj.Service.*;
import com.coxautodev.graphql.tools.GraphQLMutationResolver;
import graphql.schema.DataFetchingEnvironment;
import graphql.servlet.context.DefaultGraphQLServletContext;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import javax.servlet.http.HttpSession;

@Component
public class AdminMutation implements GraphQLMutationResolver {
    private final BulletinService bulletinService;
    private final ChallengeService challengeService;

    public AdminMutation(BulletinService bulletinService, ChallengeService challengeService) {
        this.bulletinService = bulletinService;
        this.challengeService = challengeService;
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
}

