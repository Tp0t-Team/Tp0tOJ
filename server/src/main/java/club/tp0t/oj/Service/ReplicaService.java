package club.tp0t.oj.Service;

import club.tp0t.oj.Dao.ReplicaRepository;
import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.Replica;
import club.tp0t.oj.Util.ChallengeConfiguration;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.sql.Timestamp;
import java.util.ArrayList;
import java.util.List;
import java.util.Random;

@Service
public class ReplicaService {
    @Autowired
    private ReplicaRepository replicaRepository;

    public List<Replica> getReplicaByChallenge(Challenge challenge) {
        return replicaRepository.findByChallenge(challenge);
    }

    public List<Replica> createReplicas(Challenge challenge, int count) {
        ArrayList<Replica> res = new ArrayList<>();
        for (int i = 0; i < count; i++) {
            Replica replica = new Replica();
            replica.setChallenge(challenge);
            replica.setFlag(ChallengeConfiguration.parseConfiguration(challenge.getConfiguration()).getFlag().getValue());
            res.add(replicaRepository.save(replica));
            // TODO: use saveAll to speed up
        }
        return res;
    }

    public Replica getRandomReplicaByChallenge(Challenge challenge) {
        List<Replica> replicas = replicaRepository.findAllByChallenge(challenge);
        if (replicas == null) return null;
        return replicas.get(new Random().nextInt(replicas.size()));
    }

    public void updateReplicaFlag(Replica replica, String flag) {
        replica.setFlag(flag);
        replicaRepository.save(replica);
    }
}
