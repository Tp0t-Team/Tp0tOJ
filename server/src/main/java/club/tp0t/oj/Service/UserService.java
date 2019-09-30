package club.tp0t.oj.Service;

import club.tp0t.oj.Dao.UserRepository;
import club.tp0t.oj.Entity.User;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.sql.Timestamp;
import java.util.List;
import java.util.concurrent.TimeUnit;

@Service
public class UserService {
    @Autowired
    private UserRepository userRepository;

    public List<User> getAllUsers() {
        return userRepository.findAll();
    }

    public List<User> getNormalUsers() {
        return userRepository.getNormalUsers();
    }

    public boolean checkNameExistence(String name) {
        User user = userRepository.getUserByName(name);
        return user != null;
    }

    public boolean checkStuNumberExistence(String stuNumber) {
        User user = userRepository.getUserByStuNumber(stuNumber);
        return user != null;
    }

    public boolean checkQqExistence(String qq) {
        User user = userRepository.getUserByQq(qq);
        return user != null;
    }

    public boolean checkMailExistence(String mail) {
        User user = userRepository.getUserByMail(mail);
        return user != null;
    }

    public boolean register(String name, String stuNumber, String password, String department, String qq, String mail) {
        User user = new User();
        user.setName(name);
        user.setStuNumber(stuNumber);
        user.setDepartment(department);
        user.setGmtCreated(new Timestamp(System.currentTimeMillis()));
        user.setGmtModified(new Timestamp(System.currentTimeMillis()));
        user.setJoinTime(new Timestamp(System.currentTimeMillis()));
        Timestamp timestamp = new Timestamp(System.currentTimeMillis());
        // set protected time 100 days
        timestamp.setTime(timestamp.getTime() + TimeUnit.MINUTES.toMillis(100*24*60));
        user.setProtectedTime(timestamp);
        user.setMail(mail);
        user.setPassword(password);
        user.setQQ(qq);
        user.setRole("member");
        user.setScore(0);
        user.setState("protected");
        user.setTopRank(0);

        userRepository.save(user);
        return true;
    }

    public boolean login(String name, String password) {
        User user = userRepository.getUserByName(name);

        // user disabled
        if(user.getState().equals("disabled")) {
            return false;
        }
        // password matches
        return password.equals(user.getPassword());
    }

    public boolean adminCheck(String name) {
        User user = userRepository.getUserByName(name);
        return user.getRole().equals("admin");
    }

    public long getIdByName(String name) {
        User user = userRepository.getUserByName(name);
        return user.getUserId();
    }
}
