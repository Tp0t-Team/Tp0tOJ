package club.tp0t.oj.Service;

import club.tp0t.oj.Component.ReplicaAllocHelper;
import club.tp0t.oj.Component.ReplicaHelper;
import club.tp0t.oj.Dao.ChallengeRepository;
import club.tp0t.oj.Dao.ResetTokenRepository;
import club.tp0t.oj.Dao.UserRepository;
import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.Replica;
import club.tp0t.oj.Entity.ResetToken;
import club.tp0t.oj.Entity.User;
import club.tp0t.oj.Util.FlagProxyHelper;
import club.tp0t.oj.Util.RankHelper;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Isolation;
import org.springframework.transaction.annotation.Transactional;

import java.sql.Timestamp;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;
import java.util.concurrent.TimeUnit;

@Service
public class UserService {
    private final UserRepository userRepository;
    private final ChallengeRepository challengeRepository;
    private final ResetTokenRepository resetTokenRepository;
    private final ReplicaHelper replicaHelper;
    private final ReplicaAllocHelper replicaAllocHelper;
    private final RankHelper rankHelper;
    private final FlagProxyHelper flagProxyHelper;

    public UserService(UserRepository userRepository, ChallengeRepository challengeRepository, ResetTokenRepository resetTokenRepository, ReplicaHelper replicaHelper, ReplicaAllocHelper replicaAllocHelper, RankHelper rankHelper, FlagProxyHelper flagProxyHelper) {
        this.userRepository = userRepository;
        this.challengeRepository = challengeRepository;
        this.resetTokenRepository = resetTokenRepository;
        this.replicaHelper = replicaHelper;
        this.replicaAllocHelper = replicaAllocHelper;
        this.rankHelper = rankHelper;
        this.flagProxyHelper = flagProxyHelper;
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

        rankHelper.addUser(user.getUserId(), 0);

        List<Replica> replicas = new ArrayList<>();
        for (Challenge challenge : challengeRepository.findAll()) {
            Replica replica = replicaHelper.getRandomReplicaByChallenge(challenge);
            if (replica == null) continue;
            replicas.add(replica);
        }
        //replicaAllocHelper.allocReplicasForUser(replicas, user);

        flagProxyHelper.addUser(user);
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

    public class Forget {
        private String message;
        private String token;
        private String mail;
        private String name;

        Forget(String message, String token, String mail, String name) {
            this.message = message;
            this.token = token;
            this.mail = mail;
            this.name = name;
        }

        public String getToken() {
            return token;
        }

        public String getMessage() {
            return message;
        }

        public String getMail() {
            return mail;
        }

        public String getName() {
            return name;
        }
    }

    private static String makeToken() {
        return UUID.randomUUID().toString() + "-" + Long.toString(System.currentTimeMillis(), 16);
    }

    @Transactional(isolation = Isolation.SERIALIZABLE) // for first forget
    public Forget forget(String stuNumber) {
        User user = userRepository.findByStuNumber(stuNumber);
        if (user == null) return new Forget("unregister", "", "", "");
        ResetToken resetToken = resetTokenRepository.findByUser(user);
        if (resetToken != null) {
            Timestamp now = new Timestamp(System.currentTimeMillis());
            if (now.getTime() - resetToken.getGmtModified().getTime() < TimeUnit.MINUTES.toMillis(5)) {
                // too short delta time
                return new Forget("too short interval", "", "", "");
            } else {
                // renew token
                resetToken.setToken(makeToken());
                resetToken = resetTokenRepository.save(resetToken);
                return new Forget("", resetToken.getToken(), user.getMail(), user.getName());
            }
        } else {
            // create token
            resetToken = new ResetToken();
            resetToken.setUser(user);
            resetToken.setToken(makeToken());
            resetToken = resetTokenRepository.save(resetToken);
            return new Forget("", resetToken.getToken(), user.getMail(), user.getName());
        }
    }

    @Transactional(isolation = Isolation.SERIALIZABLE) // for delete and so on
    public String reset(String token, String password) {
        ResetToken resetToken = resetTokenRepository.findByToken(token);
        if (resetToken == null) return "invalid";
        User user = resetToken.getUser();
        user.setPassword(password);
        userRepository.save(user);
        resetTokenRepository.delete(resetToken);
        return "";
    }

    public List<User> getUsersRank() {
        List<Long> rank = rankHelper.getRank();
        List<User> result = new ArrayList<>();
        for (int i = 0; i < rank.size(); i += 2) {
            long userId = rank.get(i);
            long score = rank.get(i + 1);
            User user = userRepository.findByUserId(userId);
            if (user != null) {
                user.setScore(score);
                result.add(user);
            }
        }
        return result;
        // return userRepository.getUsersRank();
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

    /*// TODO: this is an utility function.
    public void addScore(User user, long currentPoints) {
        long score = user.getScore();
        user.setScore(score + currentPoints);
        userRepository.save(user);
    }*/

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
