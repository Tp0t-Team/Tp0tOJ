package club.tp0t.oj.Dao;

import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.Replica;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface ReplicaRepository extends JpaRepository<Replica, Long> {
    List<Replica> findByChallenge(Challenge challenge);
}
