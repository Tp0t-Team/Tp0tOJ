package club.tp0t.oj.Dao;

import club.tp0t.oj.Entity.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;

import java.util.List;

public interface UserRepository extends JpaRepository<User, Long> {

    @Query(value = "select u from User u where u.role='user' and u.state='normal'")
    List<User> getNormalUsers();

    @Query(value = "select u from User u where u.name=?1")
    User getUserByName(String name);

    @Query(value = "select u from User u where u.stuNumber=?1")
    User getUserByStuNumber(String stuNumber);

    @Query(value = "select u from User u where u.QQ=?1")
    User getUserByQq(String qq);

    @Query(value = "select u from User u where u.mail=?1")
    User getUserByMail(String mail);

    @Query(value = "select u from User u where " +
            "(u.state='normal' or u.state='protected') and u.role='member' order by u.score desc ")
    List<User> getUsersRank();

    @Query(value = "select u.userId from User u where u.stuNumber=?1")
    long getUserIdByStuNumber(String stuNumber);
}
