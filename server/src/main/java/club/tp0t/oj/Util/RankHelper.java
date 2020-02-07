package club.tp0t.oj.Util;

import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.List;

@Component
public class RankHelper {
    private StringRedisTemplate redisTemplate;

    public RankHelper(StringRedisTemplate redisTemplate) {
        this.redisTemplate = redisTemplate;
    }

    public interface ScoreCalculator {
        long getScore(long challengeId, long count);

        // used for blood
        long getDeltaScoreForUser(long oldScore, long newScore, int index);
    }

    public List<Long> getRank() {
        redisTemplate.setEnableTransactionSupport(true);
        try {
            redisTemplate.multi();
            Long size = redisTemplate.opsForList().size("Rank");
            if (size == null) {
                throw new Exception();
            }
            List<String> rankList = redisTemplate.opsForList().range("Rank", 0, size);
            if (rankList == null) {
                throw new Exception();
            }
            List<Long> result = new ArrayList<>();
            for (String user : rankList) {
                result.add(Long.parseLong(user));
            }
            redisTemplate.exec();
            return result;
        } catch (Exception e) {
            redisTemplate.discard();
            return null;
        }
    }

    public boolean addUser(long userId, long baseScore) {
        redisTemplate.setEnableTransactionSupport(true);
        try {
            redisTemplate.multi();
            if (redisTemplate.hasKey("UserScore:" + userId) == Boolean.TRUE) {
                redisTemplate.opsForValue().set("UserScore:" + userId, String.valueOf(baseScore));
            }
            if (redisTemplate.hasKey("UserTime:" + userId) == Boolean.TRUE) {
                redisTemplate.opsForValue().set("UserTime:" + userId, String.valueOf(0));
            }
            redisTemplate.exec();
        } catch (Exception e) {
            redisTemplate.discard();
            return false;
        }
        return true;
    }

    public boolean addChallenge(long challenge, long originScore) {
        try {
            redisTemplate.opsForValue().set("ChallengeScore:" + challenge, String.valueOf(originScore));
        } catch (Exception e) {
            return false;
        }
        return true;
    }

    public boolean submit(long userId, long challengeId, long timestamp, ScoreCalculator calculator) {
        redisTemplate.setEnableTransactionSupport(true);
        try {
            redisTemplate.multi();
            redisTemplate.opsForValue().set("UserTime:" + userId, String.valueOf(timestamp));
            redisTemplate.opsForList().rightPush("Challenge:" + challengeId, String.valueOf(userId));
            redisTemplate.opsForList().rightPush("User:" + userId, String.valueOf(challengeId));
            String scoreString = redisTemplate.opsForValue().get("ChallengeScore:" + challengeId);
            if (scoreString == null) {
                throw new Exception();
            }
            long oldScore = Long.parseLong(scoreString);
            redisTemplate.opsForValue().increment("UserScore:" + userId, oldScore);
            Long count = redisTemplate.opsForList().size("Challenge:" + challengeId);
            if (count == null) {
                throw new Exception();
            }
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
            redisTemplate.discard();
            return false;
        }
        return true;
    }
}
