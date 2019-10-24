package club.tp0t.oj.Component;

import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.Replica;
import club.tp0t.oj.Entity.ReplicaAlloc;
import club.tp0t.oj.Entity.User;
import org.springframework.stereotype.Component;

import java.util.List;

@Component
public class FlagHelper {
    private final ReplicaAllocHelper replicaAllocHelper;
    private final ReplicaHelper replicaHelper;

    public FlagHelper(ReplicaAllocHelper replicaAllocHelper, ReplicaHelper replicaHelper) {
        this.replicaAllocHelper = replicaAllocHelper;
        this.replicaHelper = replicaHelper;
    }

    public String getFlagByUserIdAndChallengeId(User user, Challenge challenge) {
        // get user's replicas
        List<ReplicaAlloc> replicaAllocList = replicaAllocHelper.getReplicaAllocByUser(user);
        // get challenge's replicas
        List<Replica> replicas = replicaHelper.getReplicaByChallenge(challenge);

        if (replicaAllocList.size() == 0 || replicas == null) return null;

        for (ReplicaAlloc tmpReplicaAlloc : replicaAllocList) {
            for (Replica tmpReplica : replicas) {
                if (tmpReplicaAlloc.getReplica().getReplicaId() == tmpReplica.getReplicaId()) {
                    return tmpReplica.getFlag();
                }
            }
        }
        return null;
    }
}
