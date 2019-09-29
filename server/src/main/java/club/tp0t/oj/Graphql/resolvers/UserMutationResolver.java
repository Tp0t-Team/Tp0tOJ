package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Dao.*;
import club.tp0t.oj.Graphql.types.RegisterInput;
import club.tp0t.oj.Graphql.types.RegisterPayload;
import com.coxautodev.graphql.tools.GraphQLMutationResolver;
import org.springframework.stereotype.Component;

@Component
public class UserMutationResolver implements GraphQLMutationResolver {
    private final BulletinRepository bulletinRepository;
    private final ChallengeRepository challengeRepository;
    private final FlagRepository flagRepository;
    private final ReplicaRepository replicaRepository;
    private final ReplicaAllocRepository replicaAllocRepository;
    private final SubmitRepository submitRepository;
    private final UserRepository userRepository;

    public UserMutationResolver(BulletinRepository bulletinRepository,
                                ChallengeRepository challengeRepository,
                                FlagRepository flagRepository,
                                ReplicaRepository replicaRepository,
                                ReplicaAllocRepository replicaAllocRepository,
                                SubmitRepository submitRepository,
                                UserRepository userRepository) {
        this.bulletinRepository = bulletinRepository;
        this.challengeRepository = challengeRepository;
        this.flagRepository = flagRepository;
        this.replicaAllocRepository = replicaAllocRepository;
        this.replicaRepository = replicaRepository;
        this.submitRepository = submitRepository;
        this.userRepository = userRepository;
    }

    // user register
    public RegisterPayload register(RegisterInput registerInput) {

        // TODO: validate user info

        // TODO: check duplicated user

        // TODO: register user

        // if succeeded
        RegisterPayload response = new RegisterPayload("success");
        return response;
    }

}
