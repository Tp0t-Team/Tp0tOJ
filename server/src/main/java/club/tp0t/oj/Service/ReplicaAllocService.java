package club.tp0t.oj.Service;

import club.tp0t.oj.Dao.ReplicaAllocRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class ReplicaAllocService {
    @Autowired
    private ReplicaAllocRepository replicaAllocRepository;
}
