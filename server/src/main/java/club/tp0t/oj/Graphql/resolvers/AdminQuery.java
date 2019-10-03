package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Graphql.types.ChallengeConfig;
import club.tp0t.oj.Graphql.types.ChallengeConfigsResult;
import club.tp0t.oj.Service.*;
import club.tp0t.oj.Util.ChallengeConfiguration;
import com.coxautodev.graphql.tools.GraphQLQueryResolver;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.List;

@Component
public class AdminQuery implements GraphQLQueryResolver {
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

    public ChallengeConfigsResult challengeConfigs(){

        List<ChallengeConfig> challengeconfigs = new ArrayList<>();
        ChallengeConfigsResult res = new ChallengeConfigsResult("");
        List<Challenge> challenges = challengeService.getAllChallenges();

        for (Challenge challenge:challenges){
            ChallengeConfiguration challengeconfiguration  = challengeService.getConfiguration(challenge);
            ChallengeConfig challengeconfig = new ChallengeConfig();

            challengeconfig.setChallengeId(Long.toString(challenge.getChallengeId()));
            challengeconfig.setName(challengeconfiguration.getName());
            challengeconfig.setType(challengeconfiguration.getType());
            challengeconfig.setDescription(challengeconfiguration.getDescription());
            challengeconfig.setExternal_link(challengeconfiguration.getExternalLink());
            challengeconfig.setHint(challengeconfiguration.getHint());
            challengeconfig.setFlag(challengeconfiguration.getFlagEx());
            challengeconfig.setScore(challengeconfiguration.getScoreEx());

            challengeconfigs.add(challengeconfig);
        }
//        res.setChallengeConfigs(challengeconfigs);

        return res;
    }
}
