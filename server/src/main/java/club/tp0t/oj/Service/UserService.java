package club.tp0t.oj.Service;

import club.tp0t.oj.Dao.UserRepository;
import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.Submit;
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

    @Autowired
    private SubmitService submitService;

    @Autowired
    private ChallengeService challengeService;

    /*
    public List<User> getAllUsers() {
        return userRepository.findAll();
    }
    */

    /*
    public List<User> getNormalUsers() {
        return userRepository.getNormalUsers();
    }
    */

    public boolean checkNameExistence(String name) {
        //User user = userRepository.getUserByName(name);
        User user = userRepository.findByName(name);
        return user != null;
    }

    public boolean checkStuNumberExistence(String stuNumber) {
        //User user = userRepository.getUserByStuNumber(stuNumber);
        User user = userRepository.findByStuNumber(stuNumber);
        return user != null;
    }

    public boolean checkQqExistence(String qq) {
        //User user = userRepository.getUserByQq(qq);
        User user = userRepository.findByQq(qq);
        return user != null;
    }

    public boolean checkMailExistence(String mail) {
        //User user = userRepository.getUserByMail(mail);
        User user = userRepository.findByMail(mail);
        return user != null;
    }

    public boolean register(String name,
                            String stuNumber,
                            String password,
                            String department,
                            String qq,
                            String mail,
                            String grade) {
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
        user.setQq(qq);
        user.setRole("member");
        user.setScore(0);
        user.setState("protected");
        user.setTopRank(0);
        user.setGrade(grade);

        userRepository.save(user);
        return true;
    }

    public boolean login(String stuNumber, String password) {
        //User user = userRepository.getUserByStuNumber(stuNumber);
        User user = userRepository.findByStuNumber(stuNumber);

        // user disabled
        if(user.getState().equals("disabled")) {
            return false;
        }
        // password matches
        return password.equals(user.getPassword());
    }

    public boolean adminCheckByStuNumber(String stuNumber) {
        //User user = userRepository.getUserByStuNumber(stuNumber);
        User user = userRepository.findByStuNumber(stuNumber);
        return user.getRole().equals("admin");
    }

    public long getIdByName(String name) {
        //User user = userRepository.getUserByName(name);
        User user = userRepository.findByName(name);
        return user.getUserId();
    }

    public List<User> getUsersRank() {
        return userRepository.getUsersRank();
    }

    public long getIdByStuNumber(String stuNumber) {
        User user = userRepository.findByStuNumber(stuNumber);
        return user.getUserId();
    }


    public User getUserById(long userId) {
        return userRepository.getOne(userId);
    }

    public String getRoleByStuNumber(String stuNumber) {
        //User user = userRepository.getUserByStuNumber(stuNumber);
        User user = userRepository.findByStuNumber(stuNumber);
        return user.getRole();
    }

    public int getRankByUserId(long userId) {
        List<User> usersRank = userRepository.getUsersRank();
        for(int i=0;i<usersRank.size();i++) {
            User tmpUser = usersRank.get(i);
            if(tmpUser.getUserId() == userId) return i+1;
        }
        // user not exists
        return 0;
    }

    public boolean teamCheckByStuNumber(String stuNumber) {
        //User user = userRepository.getUserByStuNumber(stuNumber);
        User user = userRepository.findByStuNumber(stuNumber);
        return user.getRole().equals("team");
    }

    public boolean adminCheckByUserId(long userId) {
        User user = userRepository.getOne(userId);
        return user.getRole().equals("admin");
    }

    public boolean teamCheckByUserId(long userId) {
        User user = userRepository.getOne(userId);
        return user.getRole().equals("team");
    }

    public boolean checkUserIdExistence(long userId) {
        int count = userRepository.checkUserIdExistence(userId);
        return count == 1;
    }

    public void addScore(long userId, long currentPoints) {
        User user = userRepository.getOne(userId);
        long score = user.getScore();
        user.setScore(score + currentPoints);
        userRepository.save(user);
    }

    public void updateScore(long challengeId, long currentPoints, long newPoints) {
        Challenge challenge = challengeService.getChallengeByChallengeId(challengeId);
        List<Submit> submits = submitService.getCorrectSubmitsByChallenge(challenge);
        for (Submit tmpSubmit : submits) {
            User tmpUser = tmpSubmit.getUser();
            long score = tmpUser.getScore();
            tmpUser.setScore(score - currentPoints + newPoints);
            userRepository.save(tmpUser);
        }
    }

    public void updateUserInfo(String userId,
                               String name,
                               String role,
                               String department,
                               String grade,
                               String protectedTime,
                               String qq,
                               String mail,
                               String state) {
        User user = userRepository.getOne(Long.parseLong(userId));
        user.setName(name);
        user.setRole(role);
        user.setDepartment(department);
        user.setGrade(grade);
        user.setProtectedTime(Timestamp.valueOf(protectedTime));
        user.setQq(qq);
        user.setMail(mail);
        user.setState(state);
        userRepository.save(user);
    }
}
