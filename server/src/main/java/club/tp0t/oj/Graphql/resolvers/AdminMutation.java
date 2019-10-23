package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.Replica;
import club.tp0t.oj.Graphql.types.*;
import club.tp0t.oj.Service.*;
import com.coxautodev.graphql.tools.GraphQLMutationResolver;
import graphql.schema.DataFetchingEnvironment;
import graphql.servlet.context.DefaultGraphQLServletContext;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import javax.servlet.http.HttpSession;
import java.util.List;

@Component
public class AdminMutation implements GraphQLMutationResolver {
    private final BulletinService bulletinService;
    private final ChallengeService challengeService;
    private final FlagService flagService;
    private final ReplicaService replicaService;
    private final ReplicaAllocService replicaAllocService;
    private final SubmitService submitService;
    private final UserService userService;

    @Autowired
    public AdminMutation(BulletinService bulletinService, ChallengeService challengeService, FlagService flagService, ReplicaService replicaService, ReplicaAllocService replicaAllocService, SubmitService submitService, UserService userService) {
        this.bulletinService = bulletinService;
        this.challengeService = challengeService;
        this.flagService = flagService;
        this.replicaService = replicaService;
        this.replicaAllocService = replicaAllocService;
        this.submitService = submitService;
        this.userService = userService;
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


        if (!challengeMutate.checkPass()) return new ChallengeMutateResult("Challenge format not avaliable");

        String id = challengeMutate.getChallengeId();
        if (!id.equals("") && challengeService.checkIdExistence(Long.parseLong(id))) {
            if (!challengeService.updateChallenge(challengeMutate)) return new ChallengeMutateResult("Updation Error");
            return new ChallengeMutateResult("");
        } else {
            Challenge challenge = challengeService.addChallenge(challengeMutate);
            if (challenge == null) return new ChallengeMutateResult("Addition Error");
            List<Replica> replicas = replicaService.createReplicas(challenge, 1);
            replicaAllocService.allocReplicasForAll(replicas);
            return new ChallengeMutateResult("");
        }
    }

    public ChallengeRemoveResult challengeRemove(String id) {
        return new ChallengeRemoveResult("Can't remove any challenge");
        // TODO:
//        if(!challengeService.removeById(id)) return new ChallengeRemoveResult("Remove failed");
//        return new ChallengeRemoveResult("");
    }
}

