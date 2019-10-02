package club.tp0t.oj.Dao;

import club.tp0t.oj.Entity.Submit;
import club.tp0t.oj.Entity.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

public interface SubmitRepository extends JpaRepository<Submit, Long> {
    // TODO
    @Query(value = "select s from Submit s where s.user=?1")
    Submit getSubmitByUserIdAndChallengeId(User user, long challengeId);

    @Query(value = "select s from Submit s Join s.user su Join s.challenge sc " +
            "where su.userId=:userId and sc.challengeId=:challengeId and s.correct=true ")
    Submit checkDoneByUserId(@Param("userId")long userId,@Param("challengeId") long challengeId);
}
