package club.tp0t.oj.Graphql.types;

import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.Submit;
import club.tp0t.oj.Util.ChallengeConfiguration;
import com.alibaba.fastjson.JSON;

import java.util.ArrayList;
import java.util.List;

public class SubmitHistoryResult {
    private String message;
    private List<SubmitInfo> submitInfos;

    public SubmitHistoryResult(String message) {
        this.message = message;
        this.submitInfos = new ArrayList<>();
    }

    public void addSubmitInfos(List<Submit> submits) {
        for (Submit tmpSubmit : submits) {
            SubmitInfo tmpSubmitInfo = new SubmitInfo();
            Challenge tmpChallenge = tmpSubmit.getChallenge();

            // parse from json
            String description =  tmpChallenge.getConfiguration();
            ChallengeConfiguration challengeConfiguration = JSON.parseObject(description, ChallengeConfiguration.class);

            tmpSubmitInfo.setChallengeName(challengeConfiguration.getName());
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
