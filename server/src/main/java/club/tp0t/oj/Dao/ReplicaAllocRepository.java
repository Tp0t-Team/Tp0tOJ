package club.tp0t.oj.Dao;

import club.tp0t.oj.Entity.ReplicaAlloc;
import club.tp0t.oj.Entity.User;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface ReplicaAllocRepository extends JpaRepository<ReplicaAlloc, Long> {
    List<ReplicaAlloc> findByUser(User user);
}
