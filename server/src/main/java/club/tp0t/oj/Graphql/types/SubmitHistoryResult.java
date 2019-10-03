package club.tp0t.oj.Graphql.types;

import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.Submit;
import club.tp0t.oj.Util.ChallengeDescription;
import com.alibaba.fastjson.JSON;

import java.util.ArrayList;
import java.util.List;

public class SubmitHistoryResult {
    private String message;
    private List<SubmitInfo> submitInfos = new ArrayList<>();

    public SubmitHistoryResult(String message) {
        this.message = message;
    }

    public void addSubmitInfos(List<Submit> submits) {
        for(int i=0;i<submits.size();i++) {
            Submit tmpSubmit = submits.get(i);
            SubmitInfo tmpSubmitInfo = new SubmitInfo();
            Challenge tmpChallenge = tmpSubmit.getChallenge();

            // parse from json
            String description = tmpChallenge.getDescription();
            ChallengeDescription challengeDescription = JSON.parseObject(description, ChallengeDescription.class);

            tmpSubmitInfo.setChallengeName(challengeDescription.getName());
            tmpSubmitInfo.setMark(tmpSubmit.getMark());
            tmpSubmitInfo.setSubmitTime(tmpSubmit.getSubmitTime().toString());

            submitInfos.add(tmpSubmitInfo);
        }
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public void setSubmitInfos(List<SubmitInfo> submitInfos) {
        this.submitInfos = submitInfos;
    }

    public String getMessage() {
        return message;
    }

    public List<SubmitInfo> getSubmitInfos() {
        return submitInfos;
    }
}
