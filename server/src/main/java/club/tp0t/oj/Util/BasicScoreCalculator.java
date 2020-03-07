package club.tp0t.oj.Util;

import org.springframework.stereotype.Component;

@Component
public class BasicScoreCalculator implements RankHelper.ScoreCalculator {
    private final OjConfig ojConfig;

    public BasicScoreCalculator(OjConfig ojConfig) {
        this.ojConfig = ojConfig;
    }

    private long curve() {
        return 0;
    }

    @Override
    public long getScore(long challengeId, long count) {
        // step 1: from challengeId, get curve parameters
        // step 2: use curve(), parameters & count, calc the score
        return 0;
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
