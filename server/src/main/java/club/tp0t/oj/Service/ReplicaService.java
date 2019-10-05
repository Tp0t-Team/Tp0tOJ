package club.tp0t.oj.Service;

import club.tp0t.oj.Dao.ReplicaRepository;
import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.Replica;
import club.tp0t.oj.Entity.ReplicaAlloc;
import club.tp0t.oj.Entity.User;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class ReplicaService {
    @Autowired
    private ReplicaRepository replicaRepository;
    @Autowired
    private UserService userService;
    @Autowired
    private ChallengeService challengeService;
    @Autowired
    private ReplicaAllocService replicaAllocService;
    @Autowired
    private ReplicaService replicaService;

    public List<Replica> getReplicaByChallenge(Challenge challenge) {
        return replicaRepository.findByChallenge(challenge);
    }

    public String getFlagByUserIdAndChallengeId(long userId, long challengeId) {
        // get user entity by userId
        User user = userService.getUserById(userId);
        // user's replicas
        List<ReplicaAlloc> replicaAllocs = replicaAllocService.getReplicaAllocByUser(user);
        // get challenge entity by challengeId
        Challenge challenge = challengeService.getChallengeByChallengeId(challengeId);
        // challenge's replicas
        List<Replica> replicas = replicaService.getReplicaByChallenge(challenge);

        Replica matchReplica = null;
        for(int i=0;i<replicaAllocs.size();i++) {
            ReplicaAlloc tmpReplicaAlloc = replicaAllocs.get(i);
            for(int j=0;j<replicas.size();j++) {
                Replica tmpReplica = replicas.get(j);
                if(tmpReplicaAlloc.getReplica().getReplicaId() == tmpReplica.getReplicaId()) {
                    matchReplica = tmpReplica;
                    break;
                }
            }
            if(matchReplica != null) break;
        }
        return matchReplica.getFlag();

    }
}
