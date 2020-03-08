package club.tp0t.oj.Util;

import club.tp0t.oj.Dao.ChallengeRepository;
import club.tp0t.oj.Entity.Challenge;
import org.springframework.stereotype.Component;

@Component
public class BasicScoreCalculator implements RankHelper.ScoreCalculator {
    private final OjConfig ojConfig;
    private final ChallengeRepository challengeRepository;

    public BasicScoreCalculator(OjConfig ojConfig, ChallengeRepository challengeRepository) {
        this.ojConfig = ojConfig;
        this.challengeRepository = challengeRepository;
    }

    private long curve(int baseScore, int count) {
        if(count <= 0) {
            return baseScore;
        }
        else {
            count -= 1;
            double coefficient = 1.8414 / new Integer(ojConfig.getHalfLife()).doubleValue() * count;
            double result = Math.floor(new Integer(baseScore).doubleValue() / (coefficient + Math.exp(-coefficient)));
            return new Double(result).longValue();
        }

    }

    @Override
    public long getScore(long challengeId, long count) {
        // step 1: from challengeId, get curve parameters
        Challenge challenge = challengeRepository.findByChallengeId(challengeId);
        if(challenge == null) {
            return 0;
        }
        ChallengeConfiguration configuration = ChallengeConfiguration.parseConfiguration(challenge.getConfiguration());
        // step 2: use curve(), parameters & count, calc the score
        if(configuration.getScore().isDynamic()) {
            return curve(configuration.getScore().getBaseScore(), new Long(count).intValue());
        }
        else {
            return configuration.getScore().getBaseScore();
        }
    }

    @Override
    public long getIncrementScore(long score, long index) {
        switch (new Long(index).intValue()) {
            case 0:
                return new Double(Math.floor(score * (1 + ojConfig.getFirstBloodPercentage()))).longValue();
            case 1:
                return new Double(Math.floor(score * (1 + ojConfig.getSecondBloodPercentage()))).longValue();
            case 2:
                return new Double(Math.floor(score * (1 + ojConfig.getThirdBloodPercentage()))).longValue();
            default:
                return score;
        }
    }

    @Override
    public long getDeltaScoreForUser(long oldScore, long newScore, int index) {
        switch (index) {
            case 0:
                return (new Double(Math.floor(oldScore * (1 + ojConfig.getFirstBloodPercentage()))).longValue() - new Double(Math.floor(newScore * (1 + ojConfig.getFirstBloodPercentage()))).longValue());
            case 1:
                return (new Double(Math.floor(oldScore * (1 + ojConfig.getSecondBloodPercentage()))).longValue() - new Double(Math.floor(newScore * (1 + ojConfig.getSecondBloodPercentage()))).longValue());
            case 2:
                return (new Double(Math.floor(oldScore * (1 + ojConfig.getThirdBloodPercentage()))).longValue() - new Double(Math.floor(newScore * (1 + ojConfig.getThirdBloodPercentage()))).longValue());
            default:
                return oldScore - newScore;
        }
    }
}
