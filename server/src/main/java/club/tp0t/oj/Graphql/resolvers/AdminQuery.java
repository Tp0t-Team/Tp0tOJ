package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Service.*;
import org.springframework.beans.factory.annotation.Autowired;

public class AdminQuery extends Query {
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

}
