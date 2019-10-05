package club.tp0t.oj.Graphql.types;

import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Service.SubmitService;
import club.tp0t.oj.Util.ChallengeConfiguration;
import com.alibaba.fastjson.JSON;

import java.util.ArrayList;
import java.util.List;

public class ChallengeInfosResult {
    private String message;
    private List<ChallengeInfo> challengeInfos = new ArrayList<>();

    public ChallengeInfosResult(String message) {
        this.message = message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public void setChallengeInfos(List<ChallengeInfo> challengeInfos) {
        this.challengeInfos = challengeInfos;
    }

    public String getMessage() {
        return message;
    }

    public List<ChallengeInfo> getChallengeInfos() {
        return challengeInfos;
    }

    public void updateChallengeInfos(List<Challenge> challenges, long userId, SubmitService submitService) {
        for (int i = 0; i < challenges.size(); i++) {
            Challenge tmpChallenge = challenges.get(i);
            ChallengeInfo challengeInfo = new ChallengeInfo();
            challengeInfo.setChallengeId(Long.toString(tmpChallenge.getChallengeId()));

            // get blood
            List<String> blood = new ArrayList<>();
            if (tmpChallenge.getFirstBlood() != null)
                blood.add(Long.toString(tmpChallenge.getFirstBlood().getUserId()));
            if (tmpChallenge.getSecondBlood() != null)
                blood.add(Long.toString(tmpChallenge.getSecondBlood().getUserId()));
            if (tmpChallenge.getThirdBlood() != null)
                blood.add(Long.toString(tmpChallenge.getThirdBlood().getUserId()));
            challengeInfo.setBlood(blood);

            // whether solved by user
            challengeInfo.setDone(submitService.checkDoneByUserId(userId, tmpChallenge.getChallengeId()));

            // parse from description
            String serializedConfiguration = tmpChallenge.getConfiguration();
            ChallengeConfiguration challengeConfiguration = JSON.parseObject(serializedConfiguration, ChallengeConfiguration.class);
            challengeInfo.setDescription(challengeConfiguration.getDescription());
            challengeInfo.setExternalLink(challengeConfiguration.getExternalLink());
            challengeInfo.setHint(challengeConfiguration.getHint());
            challengeInfo.setType(challengeConfiguration.getType());
            challengeInfo.setName(challengeConfiguration.getName());
            challengeInfo.setScore(Integer.parseInt(challengeConfiguration.getScoreEx().getBase_score()));

            this.challengeInfos.add(challengeInfo);

        }
    }

}