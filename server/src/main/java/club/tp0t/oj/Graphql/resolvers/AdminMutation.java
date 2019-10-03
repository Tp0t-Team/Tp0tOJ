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
    private  SubmitService submitService;
    @Autowired
    private  UserService userService;


    public BulletinPubResult bulletinPub(BulletinPubInput bulletinPubInput, DataFetchingEnvironment environment) {

        // get session from context
        DefaultGraphQLServletContext context = environment.getContext();
        HttpSession session = context.getHttpServletRequest().getSession();

        // login & admin check
        if(session.getAttribute("isLogin") == null ||
                !((boolean) session.getAttribute("isLogin")) ||
                !(boolean) session.getAttribute("isAdmin")) {
            return new BulletinPubResult("forbidden");
        }

        String title = bulletinPubInput.getTitle();
        String content = bulletinPubInput.getContent();
        String topping = bulletinPubInput.getTopping();

        if(title==null) return new BulletinPubResult("not empty error");
        title = title.replaceAll("\\s", "");
        if(title.equals("")) return new BulletinPubResult("not empty error");
        if(bulletinService.checkTitleExistence(title)) return new BulletinPubResult("Title Already Exist");

        if(content==null) return new BulletinPubResult("not empty error");
        content = content.replaceAll("\\s", "");
        if(content.equals("")) return new BulletinPubResult("not empty error");

        if(topping==null) topping = "False";

        if(bulletinService.addBulletin(title, content, topping)) {
            return new BulletinPubResult("");
        }
        return new BulletinPubResult("Bulletin addition failed!");
    }

    public ChallengeMutateResult challengeMutate(ChallengeMutateInput challengeMutate) {
        String id = challengeMutate.getChallengeId();
        if(id==null) return new ChallengeMutateResult("not empty error");
        id  = id.replaceAll("\\s", "");
        if(id.equals("")) return new ChallengeMutateResult("not empty error");

        if(challengeService.checkIdExistence(id)){
            if(!challengeService.updateChallenge(challengeMutate)) return new ChallengeMutateResult("Updation Error");
            return new ChallengeMutateResult("");
        }else{
            if(!challengeService.checkFormat(challengeMutate)) return new ChallengeMutateResult("Challenge format not avaliable");
            if(!challengeService.addChallenge(challengeMutate)) return new ChallengeMutateResult("Addition Error");
            return new ChallengeMutateResult("");
        }

    }

    public ChallengeRemoveResult challengeRemove(String id){
        if(!challengeService.removeById(id)) return new ChallengeRemoveResult("Remove failed");
        return new ChallengeRemoveResult("");
    }
}

