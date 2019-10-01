package club.tp0t.oj.Dao;

import club.tp0t.oj.Entity.Submit;
import club.tp0t.oj.Entity.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;

public interface SubmitRepository extends JpaRepository<Submit, Long> {
    // TODO
    @Query(value = "select s from Submit s where s.user=?1 and ")
    Submit getSubmitByUserIdAndChallengeId(User user, long challengeId);
}
