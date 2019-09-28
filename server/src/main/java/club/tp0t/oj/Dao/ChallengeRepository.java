package club.tp0t.oj.Dao;

import club.tp0t.oj.Entity.Challenge;
import org.springframework.data.jpa.repository.JpaRepository;

public interface ChallengeRepository extends JpaRepository<Challenge, Long> {
}
