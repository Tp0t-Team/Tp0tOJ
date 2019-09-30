package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Graphql.types.RegisterInput;
import club.tp0t.oj.Graphql.types.RegisterPayload;
import club.tp0t.oj.Graphql.types.ResetInput;
import club.tp0t.oj.Graphql.types.ResetPayload;
import club.tp0t.oj.Service.*;
import com.coxautodev.graphql.tools.GraphQLMutationResolver;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

@Component
public class UserMutationResolver implements GraphQLMutationResolver {
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

        // TODO: validate user info

        // TODO: check duplicated user

        // TODO: register user

        // if succeeded
        return new RegisterPayload("success");
    }

    // user password reset
    public ResetPayload reset(ResetInput input) {

        // TODO: validate user info

        // TODO: reset password

        // if succeeded
        return new ResetPayload("success");

    }

}
