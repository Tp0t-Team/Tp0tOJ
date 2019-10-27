package club.tp0t.oj.Dao;

import club.tp0t.oj.Entity.ResetToken;
import club.tp0t.oj.Entity.User;
import org.springframework.data.jpa.repository.JpaRepository;

public interface ResetTokenRepository extends JpaRepository<ResetToken, Long> {
    ResetToken findByUser(User user);

    ResetToken findByToken(String token);
}
