package club.tp0t.oj.Dao;

import club.tp0t.oj.Entity.FlagProxy;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface FlagProxyRepository extends JpaRepository<FlagProxy, Long> {

    FlagProxy findByChallengeIdAndUserId(long challengeId, long userId);

    List<FlagProxy> findAllByChallengeId(Long challengeId);

    FlagProxy findByChallengeIdAndPort(Long challengeId, Long port);
}
