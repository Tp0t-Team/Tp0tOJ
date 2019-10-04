package club.tp0t.oj.Service;

import club.tp0t.oj.Dao.ReplicaAllocRepository;
import club.tp0t.oj.Entity.ReplicaAlloc;
import club.tp0t.oj.Entity.User;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class ReplicaAllocService {
    @Autowired
    private ReplicaAllocRepository replicaAllocRepository;

    public List<ReplicaAlloc> getReplicaAllocByUser(User user) {
        return replicaAllocRepository.findByUser(user);
    }
}
