package club.tp0t.oj.Util;

import club.tp0t.oj.Dao.ChallengeRepository;
import club.tp0t.oj.Dao.SubmitRepository;
import club.tp0t.oj.Dao.UserRepository;
import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.Submit;
import club.tp0t.oj.Entity.User;
import org.springframework.boot.context.event.ApplicationReadyEvent;
import org.springframework.context.event.EventListener;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.TimeUnit;

@Component
public class RankHelper {
    private final ChallengeRepository challengeRepository;
    private final UserRepository userRepository;
    private final SubmitRepository submitRepository;
    private final BasicScoreCalculator calculator;

    private StringRedisTemplate redisTemplate;

    public RankHelper(ChallengeRepository challengeRepository, UserRepository userRepository, SubmitRepository submitRepository, BasicScoreCalculator basicScoreCalculator, StringRedisTemplate redisTemplate) {
        this.challengeRepository = challengeRepository;
        this.userRepository = userRepository;
        this.submitRepository = submitRepository;
        this.calculator = basicScoreCalculator;
        this.redisTemplate = redisTemplate;
    }

    private boolean redisLock(String key, String value, long timeout) {
        if (redisTemplate.opsForValue().setIfAbsent(key, value) == Boolean.TRUE) {
            redisTemplate.expire(key, timeout, TimeUnit.MILLISECONDS);
            return true;
        }
        return false;
    }

    private void redisUnlock(String key) {
        redisTemplate.delete(key);
    }

    public interface ScoreCalculator {
        long getScore(long challengeId, long count);

        // used when this user first submit this flag
        long getIncrementScore(long score, long index); // index: 0, 1, 2, ...

        // used for blood
        long getDeltaScoreForUser(long oldScore, long newScore, int index); // index: 0, 1, 2, ...
    }

    public List<Long> getRank() {
        try {
            List<String> rankList = redisTemplate.opsForList().range("Rank", 0, -1);
            if (rankList == null) {
                throw new Exception();
            }
            List<Long> result = new ArrayList<>();
            for (String item : rankList) {
                result.add(Long.parseLong(item));
            }
            return result;
        } catch (Exception e) {
            System.out.println(e.toString());
            e.printStackTrace();
            return new ArrayList<>();
        }
    }

    public class UserInfo {
        private int rank;
        private long score;

        public int getRank() {
            return rank;
        }

        public long getScore() {
            return score;
        }

        public void setRank(int rank) {
            this.rank = rank;
        }

        public void setScore(long score) {
            this.score = score;
        }
    }

    public UserInfo getUserInfo(long userId) {
        try {
            List<String> rankList = redisTemplate.opsForList().range("Rank", 0, -1);
            if (rankList == null) {
                throw new Exception();
            }
            for (int i = 0; i < rankList.size() / 2; i++) {
                if (Long.parseLong(rankList.get(i * 2)) == userId) {
                    UserInfo info = new UserInfo();
                    info.setRank(i);
                    info.setScore(Long.parseLong(rankList.get(i * 2 + 1)));
                    return info;
                }
            }
        } catch (Exception e) {
            System.out.println(e.toString());
            e.printStackTrace();
        }
        return null;
    }

    public long getUserScore(long userId) {
        try {
            String scoreString = redisTemplate.opsForValue().get("UserScore:" + userId);
            if (scoreString == null) {
                throw new Exception();
            }
            return Long.parseLong(scoreString);
        } catch (Exception e) {
            System.out.println(e.toString());
            e.printStackTrace();
        }
        return 0;
    }

    public int getChallengeScore(long challengeId) {
        try {
            String scoreString = redisTemplate.opsForValue().get("ChallengeScore:" + challengeId);
            if (scoreString == null) {
                throw new Exception();
            }
            return Integer.parseInt(scoreString);
        } catch (Exception e) {
            System.out.println(e.toString());
            e.printStackTrace();
        }
        return 0;
    }

    public boolean addUser(long userId, long baseScore) {
        if (!redisLock("UserLock", "0", 100)) {
            return false;
        }
        try {
            redisTemplate.opsForValue().set("UserScore:" + userId, String.valueOf(baseScore));
            redisTemplate.opsForValue().set("UserTime:" + userId, String.valueOf(0));
            redisTemplate.delete("User:" + userId);
            redisUnlock("UserLock");
        } catch (Exception e) {
            e.printStackTrace();
            redisUnlock("UserLock");
            return false;
        }
        return true;
    }

    public boolean addChallenge(long challenge, long originScore) {
        if (!redisLock("ChallengeLock", "0", 100)) {
            return false;
        }
        try {
            redisTemplate.opsForValue().set("ChallengeScore:" + challenge, String.valueOf(originScore));
            redisTemplate.delete("Challenge:" + challenge);
            redisUnlock("ChallengeLock");
        } catch (Exception e) {
            e.printStackTrace();
            redisUnlock("ChallengeLock");
            return false;
        }
        return true;
    }

    public boolean submit(long userId, long challengeId, long timestamp) {
        return submit(userId, challengeId, timestamp, false);
    }

    private boolean submit(long userId, long challengeId, long timestamp, boolean warmUp) {
        if (!redisLock("SubmitLock", "0", 1000)) {
            return false;
        }
        try {
            String scoreString = redisTemplate.opsForValue().get("ChallengeScore:" + challengeId);
            if (scoreString == null) {
                throw new Exception();
            }
            Long count = redisTemplate.opsForList().size("Challenge:" + challengeId);
            if (count == null) {
                throw new Exception();
            }
            count += 1;
            long oldScore = Long.parseLong(scoreString);
            long newScore = calculator.getScore(challengeId, count);
            List<String> users = redisTemplate.opsForList().range("Challenge:" + challengeId, 0, count);
            if (users == null) {
                throw new Exception();
            }
            users.add(String.valueOf(userId));
            //
            redisTemplate.opsForValue().set("UserTime:" + userId, String.valueOf(timestamp));
            redisTemplate.opsForList().rightPush("Challenge:" + challengeId, String.valueOf(userId));
            redisTemplate.opsForList().rightPush("User:" + userId, String.valueOf(challengeId));
            redisTemplate.opsForValue().increment("UserScore:" + userId, calculator.getIncrementScore(oldScore, count - 1));
            redisTemplate.opsForValue().set("ChallengeScore:" + challengeId, String.valueOf(newScore));
            for (int i = 0; i < users.size(); i++) {
                long user = Long.parseLong(users.get(i));
                redisTemplate.opsForValue().decrement("UserScore:" + user, calculator.getDeltaScoreForUser(oldScore, newScore, i));
            }
            redisUnlock("SubmitLock");
        } catch (Exception e) {
            e.printStackTrace();
            redisUnlock("SubmitLock");
            return false;
        }
        if (warmUp) {
            return true;
        }
        return refreshRank();
    }

    private boolean refreshRank() {
        if (!redisLock("RankLock", "0", 1000)) {
            return false;
        }
        try {
            List<User> userList = userRepository.findAllByRole("member");
            List<Long[]> scoreList = new ArrayList<>();
            for (User user : userList) {
                String score = redisTemplate.opsForValue().get("UserScore:" + user.getUserId());
                if (score == null) {
                    throw new Exception();
                }
                String time = redisTemplate.opsForValue().get("UserTime:" + user.getUserId());
                if (time == null) {
                    throw new Exception();
                }
                scoreList.add(new Long[]{user.getUserId(), Long.parseLong(score), Long.parseLong(time)});
            }
            scoreList.sort((Long[] a, Long[] b) -> {
                int result = b[1].compareTo(a[1]);
                if (result == 0) {
                    result = a[2].compareTo(b[2]);
                }
                return result;
            });
            redisTemplate.delete("Rank");
            for (Long[] item : scoreList) {
                redisTemplate.opsForList().rightPush("Rank", item[0].toString());
                redisTemplate.opsForList().rightPush("Rank", item[1].toString());
            }
            redisUnlock("RankLock");
        } catch (Exception e) {
            e.printStackTrace();
            redisUnlock("RankLock");
            return false;
        }
        return true;
    }

    @EventListener(ApplicationReadyEvent.class)
    public boolean warmUp() {
        if (!redisLock("WarmUpLock", "0", 5000)) {
            return false;
        }
        List<Challenge> challengeList = challengeRepository.findByState("enabled");
        for (Challenge challenge : challengeList) {
            String challengeConfigString = challenge.getConfiguration();
            ChallengeConfiguration challengeConfiguration = ChallengeConfiguration.parseConfiguration(challengeConfigString);
            addChallenge(challenge.getChallengeId(), challengeConfiguration.getScore().getBaseScore());
        }

        List<User> userList = userRepository.findAllByRole("member");
        for (User user : userList) {
            addUser(user.getUserId(), 0/*user.getScore()*/);// TODO:
        }

        List<Submit> submitList = submitRepository.findAllByCorrectOrderBySubmitTimeAsc(true);
        for (Submit submit : submitList) {
            submit(submit.getUser().getUserId(), submit.getChallenge().getChallengeId(), submit.getSubmitTime().getTime(), true);
        }
        refreshRank();
        redisUnlock("WarmUpLock");
        return true;
    }
}
