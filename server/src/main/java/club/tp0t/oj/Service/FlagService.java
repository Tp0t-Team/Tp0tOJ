package club.tp0t.oj.Service;

import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.Replica;
import club.tp0t.oj.Entity.ReplicaAlloc;
import club.tp0t.oj.Entity.User;
import org.springframework.stereotype.Service;

import java.util.List;

// TODO: all is utility function.
@Service
public class FlagService {
    private final ReplicaAllocService replicaAllocService;
    private final ReplicaService replicaService;

    public FlagService(ReplicaAllocService replicaAllocService, ReplicaService replicaService) {
        this.replicaAllocService = replicaAllocService;
        this.replicaService = replicaService;
    }

    public String getFlagByUserIdAndChallengeId(User user, Challenge challenge) {
        //Flag flag = flagRepository.getFlagByUserIdAndChallengeId(userId, challengeId);
        //Flag flag = flagRepository.test(challengeId);

        // get user entity by userId
//        User user = userService.getUserById(userId);
        // user's replicas
        List<ReplicaAlloc> replicaAllocs = replicaAllocService.getReplicaAllocByUser(user);
        // get challenge entity by challengeId
//        Challenge challenge = challengeService.getChallengeByChallengeId(challengeId);
        // challenge's replicas
        List<Replica> replicas = replicaService.getReplicaByChallenge(challenge);

        if (replicaAllocs == null || replicas == null) return null;

        for (ReplicaAlloc tmpReplicaAlloc : replicaAllocs) {
            for (Replica tmpReplica : replicas) {
                if (tmpReplicaAlloc.getReplica().getReplicaId() == tmpReplica.getReplicaId()) {
                    return tmpReplica.getFlag();
                }
            }
        }
        return null;
    }

}
