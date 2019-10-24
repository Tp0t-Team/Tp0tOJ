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

}