package club.tp0t.oj.Service;

import club.tp0t.oj.Dao.SubmitRepository;
import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.Submit;
import club.tp0t.oj.Entity.User;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.sql.Timestamp;
import java.util.List;

@Service
public class SubmitService {
    @Autowired
    private SubmitRepository submitRepository;
    @Autowired
    private ChallengeService challengeService;
    @Autowired
    private UserService userService;

    public void submit(User user, Challenge challenge, String submitFlag, boolean correct, int mark) {
        Submit submit = new Submit();
        submit.setCorrect(correct);
        submit.setGmtCreated(new Timestamp(System.currentTimeMillis()));
        submit.setGmtModified(new Timestamp(System.currentTimeMillis()));
        submit.setMark(mark);
        submit.setSubmitFlag(submitFlag);
        submit.setSubmitTime(new Timestamp(System.currentTimeMillis()));
        submit.setUser(user);
        submit.setChallenge(challenge);

        submitRepository.save(submit);
    }

    public boolean checkDuplicateSubmit(User user, long challengeId) {
        //Submit submit = submitRepository.getSubmitByUserIdAndChallengeId(user, challengeId);
        Challenge challenge = challengeService.getChallengeByChallengeId(challengeId);
        Submit submit = submitRepository.findByUserAndChallengeAndCorrect(user, challenge, true);
        // duplicate: true
        return submit != null;
    }


    public boolean checkDoneByUserId(long userId, long challengeId) {
        User user = userService.getUserById(userId);
        Challenge challenge = challengeService.getChallengeByChallengeId(challengeId);
        //Submit submit = submitRepository.checkDoneByUserId(userId, challengeId);
        Submit submit = submitRepository.findByUserAndChallengeAndCorrect(user, challenge, true);
//        if (submit == null) return true;
//        else return false;
        return submit != null;
    }

    public List<Submit> getCorrectSubmitsByChallenge(Challenge challenge) {
        //return submitRepository.getCorrectSubmitsByChallenge(challenge);
        return submitRepository.findAllByChallengeAndCorrect(challenge, true);
    }

    @Transactional
    public int updateSolvedCountByChallengeId(long challengeId, User user) {
        Challenge challenge = challengeService.getChallengeByChallengeId(challengeId);
        //List<Submit> submits = submitRepository.getCorrectSubmitsByChallenge(challenge);
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
        return submits.size() + 1;
    }

    public List<Submit> getCorrectSubmitsByUserId(long userId) {
        User user = userService.getUserById(userId);
        if (user == null) return null;
        return submitRepository.findAllByUserAndCorrect(user, true);
    }
}
