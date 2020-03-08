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
            System.out.println(rankList.size());
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

    public boolean addUser(long userId, long baseScore) {
        if (!redisLock("UserLock", "0", 100)) {
            return false;
        }
        // redisTemplate.setEnableTransactionSupport(true);
        try {
            // redisTemplate.multi();
            redisTemplate.opsForValue().set("UserScore:" + userId, String.valueOf(baseScore));
            redisTemplate.opsForValue().set("UserTime:" + userId, String.valueOf(0));
            redisTemplate.delete("User:" + userId);
            // redisTemplate.exec();
            redisUnlock("UserLock");
        } catch (Exception e) {
            e.printStackTrace();
            // redisTemplate.discard();
            redisUnlock("UserLock");
            return false;
        }
        return true;
    }

    public boolean addChallenge(long challenge, long originScore) {
        if (!redisLock("ChallengeLock", "0", 100)) {
            return false;
        }
        // redisTemplate.setEnableTransactionSupport(true);
        try {
            // redisTemplate.multi();
            redisTemplate.opsForValue().set("ChallengeScore:" + challenge, String.valueOf(originScore));
            redisTemplate.delete("Challenge:" + challenge);
            redisUnlock("ChallengeLock");
            // redisTemplate.exec();
        } catch (Exception e) {
            e.printStackTrace();
            // redisTemplate.discard();
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
        // redisTemplate.setEnableTransactionSupport(true);
        try {
            /*while (redisTemplate.opsForValue().setIfAbsent("lock", String.valueOf(1)) == Boolean.TRUE) {
                try {
                    wait(10);
                } catch (Exception e) {
                    e.printStackTrace();
                }
            }
            redisTemplate.expire("lock", 2, TimeUnit.SECONDS);*/
            // redisTemplate.watch("ChallengeScore:" + challengeId);
            // redisTemplate.watch("Challenge:" + challengeId);
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
            try {
                // redisTemplate.multi();
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
                // redisTemplate.exec();
                redisUnlock("SubmitLock");
            } catch (Exception e) {
                e.printStackTrace();
                // redisTemplate.discard();
                redisUnlock("SubmitLock");
                return false;
            }
            // redisTemplate.unwatch();
            //redisTemplate.delete("lock");
        } catch (Exception e) {
            //redisTemplate.delete("lock");
            e.printStackTrace();
            // redisTemplate.unwatch();
            redisUnlock("SubmitLock");
            return false;
        }
        if (warmUp) {
            return true;
        }
        return refreshRank();
        /*try {
            redisTemplate.multi();
            redisTemplate.opsForValue().set("UserTime:" + userId, String.valueOf(timestamp));
            redisTemplate.opsForList().rightPush("Challenge:" + challengeId, String.valueOf(userId));
            redisTemplate.opsForList().rightPush("User:" + userId, String.valueOf(challengeId));
            String scoreString = redisTemplate.opsForValue().get("ChallengeScore:" + challengeId);
            if (scoreString == null) {
                System.out.println(challengeId);
                throw new Exception();
            }
            Long count = redisTemplate.opsForList().size("Challenge:" + challengeId);
            if (count == null) {
                throw new Exception();
            }
            long oldScore = Long.parseLong(scoreString);
            redisTemplate.opsForValue().increment("UserScore:" + userId, calculator.getIncrementScore(oldScore, count - 1));
            long newScore = calculator.getScore(challengeId, count);
            redisTemplate.opsForValue().set("ChallengeScore:" + challengeId, String.valueOf(newScore));
            List<String> users = redisTemplate.opsForList().range("Challenge:", 0, count);
            if (users == null) {
                throw new Exception();
            }
            List<Long[]> scoreList = new ArrayList<>();
            for (int i = 0; i < users.size(); i++) {
                long user = Long.parseLong(users.get(i));
                redisTemplate.opsForValue().decrement("UserScore:" + user, calculator.getDeltaScoreForUser(oldScore, newScore, i));
                String score = redisTemplate.opsForValue().get("UserScore:" + user);
                if (score == null) {
                    throw new Exception();
                }
                String time = redisTemplate.opsForValue().get("UserTime:" + user);
                if (time == null) {
                    throw new Exception();
                }
                scoreList.add(new Long[]{user, Long.parseLong(score), Long.parseLong(time)});
            }
            scoreList.sort((Long[] a, Long[] b) -> {
                int result = b[1].compareTo(a[1]);
                if (result == 0) {
                    result = a[1].compareTo(b[1]);
                }
                return result;
            });
            for (int i = 0; i < scoreList.size(); i++) {
                redisTemplate.opsForList().set("Rank", i, scoreList.get(i)[0].toString());
            }
            redisTemplate.exec();
        } catch (Exception e) {
            System.out.println(e.toString());
            e.printStackTrace();
            redisTemplate.discard();
            return false;
        }*/
    }

    private boolean refreshRank() {
        if (!redisLock("RankLock", "0", 1000)) {
            return false;
        }
        // redisTemplate.setEnableTransactionSupport(true);
        try {
            List<User> userList = userRepository.findAllByRole("member");
            List<Long[]> scoreList = new ArrayList<>();
            for (User user : userList) {
                // redisTemplate.watch("UserScore:" + user.getUserId());
                String score = redisTemplate.opsForValue().get("UserScore:" + user.getUserId());
                if (score == null) {
                    throw new Exception();
                }
                // redisTemplate.watch("UserTime:" + user.getUserId());
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
            try {
                // redisTemplate.multi();
                redisTemplate.delete("Rank");
                for (Long[] item : scoreList) {
                    redisTemplate.opsForList().rightPush("Rank", item[0].toString());
                    redisTemplate.opsForList().rightPush("Rank", item[1].toString());
                }
                redisUnlock("RankLock");
                // redisTemplate.exec();
            } catch (Exception e) {
                e.printStackTrace();
                // redisTemplate.discard();
                redisUnlock("RankLock");
                return false;
            }
            // redisTemplate.unwatch();
        } catch (Exception e) {
            e.printStackTrace();
            // redisTemplate.unwatch();
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
