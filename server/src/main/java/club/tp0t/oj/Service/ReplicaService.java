package club.tp0t.oj.Service;

import club.tp0t.oj.Dao.ReplicaRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class ReplicaService {
    @Autowired
    private ReplicaRepository replicaRepository;
}
