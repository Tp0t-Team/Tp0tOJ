package club.tp0t.oj.Service;

import club.tp0t.oj.Dao.SubmitRepository;
import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.Submit;
import club.tp0t.oj.Entity.User;
import club.tp0t.oj.Util.ChallengeConfiguration;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Isolation;
import org.springframework.transaction.annotation.Transactional;

import java.sql.Timestamp;
import java.util.List;

@Service
public class SubmitService {
    private final SubmitRepository submitRepository;
    private final ChallengeService challengeService;
    private final UserService userService;
    private final FlagService flagService;

    public SubmitService(SubmitRepository submitRepository, ChallengeService challengeService, UserService userService, FlagService flagService) {
        this.submitRepository = submitRepository;
        this.challengeService = challengeService;
        this.userService = userService;
        this.flagService = flagService;
    }

    @Transactional(isolation = Isolation.SERIALIZABLE)
    // maybe can user lower level, but I think use the safest level is ok.
    public String submit(long userId, long challengeId, String flag, boolean isMember) {
        // prepare
        User user = userService.getUserById(userId);
        Challenge challenge = challengeService.getChallengeByChallengeId(challengeId);

        // get your flag and check
        String correctFlag = flagService.getFlagByUserIdAndChallengeId(user, challenge);
        if (correctFlag == null) return "No replica for you.";
        boolean correct = correctFlag.equals(flag);

        if (!isMember) {
            if (correct) {
                return "correct";
            } else {
                return "incorrect";
            }
        } else {
            int mark = 0;
            if (correct) {
                if (submitRepository.findByUserAndChallengeAndCorrect(user, challenge, true) != null) {
                    return "duplicate submit";
                }
                // add user score
                // TODO: get current points of challenge
                long currentPoints = 0;
                try {
                    currentPoints = Long.parseLong(ChallengeConfiguration.parseConfiguration(challenge.getConfiguration()).getScoreEx().getBase_score());
                } catch (NumberFormatException e) {
                    return "unknown error";
                }
                userService.addScore(userId, currentPoints);

                // update blood
                List<Submit> submits = submitRepository.findAllByChallengeAndCorrect(challenge, true);
                switch (submits.size()) {
                    case 0:
                        challenge.setFirstBlood(user);
                        challengeService.updateChallengeBlood(challenge);
                        break;
                    case 1:
                        challenge.setSecondBlood(user);
                        challengeService.updateChallengeBlood(challenge);
                        break;
                    case 2:
                        challenge.setThirdBlood(user);
                        challengeService.updateChallengeBlood(challenge);
                        break;
                    default:
                        break;
                }
                if (submits.size() < 3) {
                    mark = submits.size() + 1;
                }
            }
            Submit submit = new Submit();
            submit.setCorrect(correct);
            submit.setMark(mark);
            submit.setSubmitFlag(flag);
            submit.setSubmitTime(new Timestamp(System.currentTimeMillis()));
            submit.setUser(user);
            submit.setChallenge(challenge);
            submitRepository.save(submit);
            if (correct) {
                return "";
            } else {
                return "incorrect";
            }
        }
    }

    public List<Submit> getCorrectSubmitsByUserId(long userId) {
        User user = userService.getUserById(userId);
        if (user == null) return null;
        return submitRepository.findAllByUserAndCorrect(user, true);
    }
}
