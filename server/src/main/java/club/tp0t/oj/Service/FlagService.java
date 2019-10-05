package club.tp0t.oj.Service;

import club.tp0t.oj.Dao.FlagRepository;
import club.tp0t.oj.Entity.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class FlagService {
    @Autowired
    private FlagRepository flagRepository;
    @Autowired
    private ReplicaAllocService replicaAllocService;
    @Autowired
    private ReplicaService replicaService;
    @Autowired
    private UserService userService;
    @Autowired
    private ChallengeService challengeService;

    public String getFlagByUserIdAndChallengeId(long userId, long challengeId) {
        //Flag flag = flagRepository.getFlagByUserIdAndChallengeId(userId, challengeId);
        //Flag flag = flagRepository.test(challengeId);

        // get user entity by userId
        User user = userService.getUserById(userId);
        // user's replicas
        List<ReplicaAlloc> replicaAllocs = replicaAllocService.getReplicaAllocByUser(user);
        // get challenge entity by challengeId
        Challenge challenge = challengeService.getChallengeByChallengeId(challengeId);
        // challenge's replicas
        List<Replica> replicas = replicaService.getReplicaByChallenge(challenge);

        Replica matchReplica = null;
        for (ReplicaAlloc tmpReplicaAlloc : replicaAllocs) {
            for (Replica tmpReplica : replicas) {
                if (tmpReplicaAlloc.getReplica().getReplicaId() == tmpReplica.getReplicaId()) {
                    matchReplica = tmpReplica;
                    break;
                }
            }
            if (matchReplica != null) break;
        }
        Flag flag = flagRepository.findByReplica(matchReplica);
        return flag.getFlag();

    }

}
