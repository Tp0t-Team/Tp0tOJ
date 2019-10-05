package club.tp0t.oj.Service;

import club.tp0t.oj.Dao.ReplicaAllocRepository;
import club.tp0t.oj.Entity.Replica;
import club.tp0t.oj.Entity.ReplicaAlloc;
import club.tp0t.oj.Entity.User;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.sql.Timestamp;
import java.util.List;

@Service
public class ReplicaAllocService {
    @Autowired
    private ReplicaAllocRepository replicaAllocRepository;
    @Autowired
    private UserService userService;

    public List<ReplicaAlloc> getReplicaAllocByUser(User user) {
        return replicaAllocRepository.findByUser(user);
    }

    public void allocReplicasForAll(List<Replica> replicas) {
        List<User> users = userService.getAllUser();
        if(users == null) return;
        int userPerReplica = replicas.size() / users.size();
        for (int i = 0; i < replicas.size(); i++) {
            for (int j = 0; j < userPerReplica; j++) {
                ReplicaAlloc replicaAlloc = new ReplicaAlloc();
                replicaAlloc.setReplica(replicas.get(i));
                replicaAlloc.setUser(users.get(i * userPerReplica + j));
                replicaAlloc.setGmtCreated(new Timestamp(System.currentTimeMillis()));
                replicaAlloc.setGmtModified(new Timestamp(System.currentTimeMillis()));
                replicaAllocRepository.save(replicaAlloc);
            }
        }
        for (int i = replicas.size() * userPerReplica; i < users.size(); i++) {
            ReplicaAlloc replicaAlloc = new ReplicaAlloc();
            replicaAlloc.setReplica(replicas.get(0));
            replicaAlloc.setUser(users.get(i));
            replicaAlloc.setGmtCreated(new Timestamp(System.currentTimeMillis()));
            replicaAlloc.setGmtModified(new Timestamp(System.currentTimeMillis()));
            replicaAllocRepository.save(replicaAlloc);
        }
    }

    public void allocReplicasForUser(List<Replica> replicas, User user) {
        for (Replica replica : replicas) {
            ReplicaAlloc replicaAlloc = new ReplicaAlloc();
            replicaAlloc.setReplica(replica);
            replicaAlloc.setUser(user);
            replicaAlloc.setGmtCreated(new Timestamp(System.currentTimeMillis()));
            replicaAlloc.setGmtModified(new Timestamp(System.currentTimeMillis()));
            replicaAllocRepository.save(replicaAlloc);
        }
    }
}
