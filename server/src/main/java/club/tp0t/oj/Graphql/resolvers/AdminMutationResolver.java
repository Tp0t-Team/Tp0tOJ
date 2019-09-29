package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Dao.*;
import com.coxautodev.graphql.tools.GraphQLMutationResolver;

public class AdminMutationResolver implements GraphQLMutationResolver {
    private final BulletinRepository bulletinRepository;
    private final ChallengeRepository challengeRepository;
    private final FlagRepository flagRepository;
    private final ReplicaRepository replicaRepository;
    private final ReplicaAllocRepository replicaAllocRepository;
    private final SubmitRepository submitRepository;
    private final UserRepository userRepository;

    public AdminMutationResolver(BulletinRepository bulletinRepository,
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
}
