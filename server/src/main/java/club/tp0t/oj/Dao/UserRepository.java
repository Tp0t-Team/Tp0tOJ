package club.tp0t.oj.Dao;

import club.tp0t.oj.Entity.User;
import org.springframework.data.jpa.repository.JpaRepository;

public interface UserRepository extends JpaRepository<User, Long> {
}
