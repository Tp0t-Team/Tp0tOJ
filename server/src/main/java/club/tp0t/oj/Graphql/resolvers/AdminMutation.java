package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Graphql.types.*;
import club.tp0t.oj.Service.*;
import com.coxautodev.graphql.tools.GraphQLMutationResolver;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

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


    public BulletinPubResult bulletinPub(BulletinPubInput bulletinPubInput) {


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
}
