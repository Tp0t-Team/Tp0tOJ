package club.tp0t.oj.Graphql.types;

import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Service.SubmitService;
import club.tp0t.oj.Util.ChallengeDescription;
import com.alibaba.fastjson.JSON;

import java.util.ArrayList;
import java.util.List;

public class ChallengesResult {
    private String message;
    private List<ChallengeInfo> challengeInfos = new  ArrayList<>();

    public ChallengesResult(String message) {
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

    public void addChallengeInfos(List<Challenge> challenges, long userId, SubmitService submitService) {
        for(int i=0;i<challenges.size();i++) {
            System.out.println(i);
            Challenge tmpChallenge = challenges.get(i);
            ChallengeInfo challengeInfo = new ChallengeInfo();
            challengeInfo.setChallengeId(Long.toString(tmpChallenge.getChallengeId()));

            // get blood
            List<String> blood = new ArrayList<>();
            if(tmpChallenge.getFirstBlood() != null) blood.add(tmpChallenge.getFirstBlood().getName());
            if(tmpChallenge.getSecondBlood() != null) blood.add(tmpChallenge.getSecondBlood().getName());
            if(tmpChallenge.getThirdBlood() != null) blood.add(tmpChallenge.getThirdBlood().getName());
            challengeInfo.setBlood(blood);

            // whether solved by user
            challengeInfo.setDone(submitService.checkDoneByUserId(userId, tmpChallenge.getChallengeId()));

            // parse from description
            String description = tmpChallenge.getDescription();
            ChallengeDescription challengeDescription = JSON.parseObject(description, ChallengeDescription.class);
            challengeInfo.setDescription(challengeDescription.getDescription());
            challengeInfo.setExternalLink(challengeDescription.getExternalLink());
            challengeInfo.setHint(challengeDescription.getHint());
            challengeInfo.setType(challengeDescription.getType());
            challengeInfo.setName(challengeDescription.getName());
            challengeInfo.setScore(challengeDescription.getScore());

            this.challengeInfos.add(challengeInfo);

        }
    }

}
