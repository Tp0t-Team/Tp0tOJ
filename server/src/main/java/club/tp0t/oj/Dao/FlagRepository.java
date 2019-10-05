//package club.tp0t.oj.Dao;
//
//import club.tp0t.oj.Entity.Replica;
//import org.springframework.data.jpa.repository.JpaRepository;
//import org.springframework.data.jpa.repository.Query;
//import org.springframework.data.repository.query.Param;
//
//public interface FlagRepository extends JpaRepository<Flag, Long> {
//
//    /*
//    @Query(value = "select f from Flag f Join f.challenge fc Join f.replica fr where " +
//            "fc.challengeId=:challengeId and " +
//            "fr.replicaId in (select ra.replicaId from ReplicaAlloc ra where ra.userId=:userId)"
//            )
//    Flag getFlagByUserIdAndChallengeId(@Param("userId") long userId, @Param("challengeId") long challengeId);
//
//    */
//
//    Flag findByReplica(Replica matchReplica);
//}
//
