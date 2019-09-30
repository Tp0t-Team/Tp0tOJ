package club.tp0t.oj.Dao;

import club.tp0t.oj.Entity.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;

import java.util.List;

public interface UserRepository extends JpaRepository<User, Long> {

    @Query(value = "select u from User u where u.role='user' and u.state='normal'")
    List<User> getNormalUsers();
}
