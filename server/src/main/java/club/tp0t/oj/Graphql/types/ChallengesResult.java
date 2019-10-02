package club.tp0t.oj.Graphql.types;

import java.util.ArrayList;
import java.util.List;

public class ChallengesResult {
    private String message;
    private List<ChallengeInfo> challengeInfos = new  ArrayList<>();

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
