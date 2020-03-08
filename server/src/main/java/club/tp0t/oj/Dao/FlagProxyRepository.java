package club.tp0t.oj.Dao;

import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.FlagProxy;
import club.tp0t.oj.Entity.User;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface FlagProxyRepository extends JpaRepository<FlagProxy, Long> {

    FlagProxy findByChallengeAndUser(Challenge challenge, User user);

    List<FlagProxy> findAllByChallenge(Challenge challenge);

    FlagProxy findByChallengeAndPort(Challenge challenge, Long port);
}
