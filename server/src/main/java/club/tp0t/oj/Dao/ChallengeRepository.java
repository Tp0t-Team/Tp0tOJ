package club.tp0t.oj.Dao;

import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Graphql.types.ChallengeInfo;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;

import javax.validation.constraints.NotEmpty;
import java.util.List;

public interface ChallengeRepository extends JpaRepository<Challenge, Long> {
    //@Query(value = "select c from Challenge c where c.state='enabled'")
    //List<Challenge> getEnabledChallenges();

    List<Challenge> findByState(String state);

    Challenge findByChallengeId(String id);
}
