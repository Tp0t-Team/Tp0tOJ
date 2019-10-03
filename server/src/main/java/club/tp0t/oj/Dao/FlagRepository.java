package club.tp0t.oj.Dao;

import club.tp0t.oj.Entity.Flag;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

public interface FlagRepository extends JpaRepository<Flag, Long> {

    /*
    @Query(value = "select f.flag from Flag f Join f.challenge fc Join f.replica fr Join ReplicaAlloc ra " +
            "where ra.replicaAllocId=fr.replicaId and fc.")
    String getFlagByUserIdAndChallengeId(@Param("userId") long userId, @Param("challengeId") long challengeId);
    */
}

