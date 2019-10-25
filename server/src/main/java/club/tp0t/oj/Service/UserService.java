package club.tp0t.oj.Service;

import club.tp0t.oj.Component.ReplicaAllocHelper;
import club.tp0t.oj.Component.ReplicaHelper;
import club.tp0t.oj.Dao.ChallengeRepository;
import club.tp0t.oj.Dao.UserRepository;
import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.Replica;
import club.tp0t.oj.Entity.User;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Isolation;
import org.springframework.transaction.annotation.Transactional;

import java.sql.Timestamp;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.TimeUnit;

@Service
public class UserService {
    private final UserRepository userRepository;
    private final ChallengeRepository challengeRepository;
    private final ReplicaHelper replicaHelper;
    private final ReplicaAllocHelper replicaAllocHelper;

    public UserService(UserRepository userRepository, ChallengeRepository challengeRepository, ReplicaHelper replicaHelper, ReplicaAllocHelper replicaAllocHelper) {
        this.userRepository = userRepository;
        this.challengeRepository = challengeRepository;
        this.replicaHelper = replicaHelper;
        this.replicaAllocHelper = replicaAllocHelper;
    }

    private boolean checkStuNumberExistence(String stuNumber) {
        //User user = userRepository.getUserByStuNumber(stuNumber);
        User user = userRepository.findByStuNumber(stuNumber);
        return user != null;
    }

    private boolean checkQqExistence(String qq) {
        //User user = userRepository.getUserByQq(qq);
        User user = userRepository.findByQq(qq);
        return user != null;
    }

    private boolean checkMailExistence(String mail) {
        //User user = userRepository.getUserByMail(mail);
        User user = userRepository.findByMail(mail);
        return user != null;
    }

    @Transactional(isolation = Isolation.SERIALIZABLE) // for unique test, must use this level.
    public String register(String name,
                           String stuNumber,
                           String password,
                           String department,
                           String qq,
                           String mail,
                           String grade) {
        if (checkStuNumberExistence(stuNumber)) {
            return "Student number has been used.";
        }
        if (checkQqExistence(qq)) {
            return "QQ has been used.";
        }
        if (checkMailExistence(mail)) {
            return "Mail has been used.";
        }
        User user = new User();
        user.setName(name);
        user.setStuNumber(stuNumber);
        user.setDepartment(department);
        user.setJoinTime(new Timestamp(System.currentTimeMillis()));
        Timestamp timestamp = new Timestamp(System.currentTimeMillis());
        // set protected time 100 days
        // TODO: set to correct protected time!!!
        timestamp.setTime(timestamp.getTime() + TimeUnit.MINUTES.toMillis(100 * 24 * 60));
        user.setProtectedTime(timestamp);
        user.setMail(mail);
        user.setPassword(password);
        user.setQq(qq);
        user.setRole("member");
        user.setScore(0);
        user.setState("protected");
        user.setTopRank(0);
        user.setGrade(grade);

        user = userRepository.save(user);

        List<Replica> replicas = new ArrayList<>();
        for (Challenge challenge : challengeRepository.findAll()) {
            Replica replica = replicaHelper.getRandomReplicaByChallenge(challenge);
            if (replica == null) continue;
            replicas.add(replica);
        }
        replicaAllocHelper.allocReplicasForUser(replicas, user);

        return "";
    }

    @Transactional
    public User login(String stuNumber, String password) {
        User user = userRepository.findByStuNumber(stuNumber);
        // not exist
        if (user == null) return null;
        // user disabled
        if (user.getState().equals("disabled")) {
            return null;
        }
        // password matches
        if (password.equals(user.getPassword())) {
            return user;
        } else {
            return null;
        }
    }

    public List<User> getUsersRank() {
        return userRepository.getUsersRank();
    }

    public List<User> getAllUser() {
        return userRepository.findAll();
    }

    // TODO: is this utility function necessary?
    public User getUserById(long userId) {
        return userRepository.findByUserId(userId);
    }

    public int getRankByUserId(long userId) {
        List<User> usersRank = userRepository.getUsersRank();
        for (int i = 0; i < usersRank.size(); i++) {
            User tmpUser = usersRank.get(i);
            if (tmpUser.getUserId() == userId) return i + 1;
        }
        // user not exists
        return 0;
    }

    // TODO: this is an utility function.
    public boolean adminCheckByUserId(long userId) {
        User user = userRepository.getOne(userId);
        return user.getRole().equals("admin");
    }

//    public boolean teamCheckByUserId(long userId) {
//        User user = userRepository.getOne(userId);
//        return user.getRole().equals("team");
//    }

    public boolean checkUserIdExistence(long userId) {
        //int count = userRepository.checkUserIdExistence(userId);
        long count = userRepository.countByUserId(userId);
        return count == 1;
    }

    // TODO: this is an utility function.
    public void addScore(User user, long currentPoints) {
        long score = user.getScore();
        user.setScore(score + currentPoints);
        userRepository.save(user);
    }

//    public void updateScore(long challengeId, long currentPoints, long newPoints) {
//        Challenge challenge = challengeRepository.findByChallengeId(challengeId);// challengeService.getChallengeByChallengeId(challengeId);
////        List<Submit> submits = submitService.getCorrectSubmitsByChallenge(challenge);
//        List<Submit> submits = submitRepository.findAllByChallengeAndCorrect(challenge, true);
//        for (Submit tmpSubmit : submits) {
//            User tmpUser = tmpSubmit.getUser();
//            long score = tmpUser.getScore();
//            tmpUser.setScore(score - currentPoints + newPoints);
//            userRepository.save(tmpUser);
//        }
//    }

    @Transactional(isolation = Isolation.REPEATABLE_READ)
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
