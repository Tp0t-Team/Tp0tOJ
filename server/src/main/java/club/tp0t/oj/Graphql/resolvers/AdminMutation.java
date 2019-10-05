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
    private SubmitService submitService;
    @Autowired
    private UserService userService;


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

        if (!bulletinPubInput.checkPass()) return new BulletinPubResult("not empty error");

        String title = bulletinPubInput.getTitle();
        String content = bulletinPubInput.getContent();
        boolean topping = bulletinPubInput.getTopping();

//        if(title==null) return new BulletinPubResult("not empty error");
//        title = title.trim();
//        if(title.equals("")) return new BulletinPubResult("not empty error");
//        if(content==null) return new BulletinPubResult("not empty error");
//        content = content.trim();
//        if(content.equals("")) return new BulletinPubResult("not empty error");

        if (bulletinService.addBulletin(title, content, topping)) {
            return new BulletinPubResult("");
        }
        return new BulletinPubResult("Bulletin addition failed!");
    }

    public ChallengeMutateResult challengeMutate(ChallengeMutateInput challengeMutate, DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // login & admin check
        if (session.getAttribute("isLogin") == null ||
                !((boolean) session.getAttribute("isLogin")) ||
                !(boolean) session.getAttribute("isAdmin")) {
            return new ChallengeMutateResult("forbidden");
        }


        if (!challengeMutate.checkPass()) return new ChallengeMutateResult("Challenge format not avaliable");
        
        String id = challengeMutate.getChallengeId();
        if (!id.equals("") && challengeService.checkIdExistence(Long.parseLong(id))) {
            if (!challengeService.updateChallenge(challengeMutate)) return new ChallengeMutateResult("Updation Error");
            return new ChallengeMutateResult("");
        } else {
            if (!challengeService.addChallenge(challengeMutate)) return new ChallengeMutateResult("Addition Error");
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

