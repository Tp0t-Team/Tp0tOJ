package club.tp0t.oj.Service;

import club.tp0t.oj.Dao.SubmitRepository;
import club.tp0t.oj.Entity.Submit;
import club.tp0t.oj.Entity.User;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.sql.Timestamp;

@Service
public class SubmitService {
    @Autowired
    private SubmitRepository submitRepository;

    public void submit(User user, String submitFlag, boolean correct, int mark) {
        Submit submit = new Submit();
        submit.setCorrect(correct);
        submit.setGmtCreated(new Timestamp(System.currentTimeMillis()));
        submit.setGmtModified(new Timestamp(System.currentTimeMillis()));
        submit.setMark(mark);
        submit.setSubmitFlag(submitFlag);
        submit.setSubmitTime(new Timestamp(System.currentTimeMillis()));
        submit.setUser(user);
    }

    public boolean checkDuplicateSubmit(User user, long challengeId) {
        Submit submit = submitRepository.getSubmitByUserIdAndChallengeId(user, challengeId);
        // duplicate: true
        return submit != null;
    }

    public boolean checkDoneByUserId(long userId, long challengeId) {
        Submit submit = submitRepository.checkDoneByUserId(userId, challengeId);
        if(submit == null) return true;
        else return false;
    }
}
