package club.tp0t.oj.Graphql.types;

import java.util.ArrayList;
import java.util.List;

public class ChallengeConfigsResult {
    private String message;
    private List<ChallengeConfig> challengeConfigs;

    public ChallengeConfigsResult(String message){
        this.message = message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public String getMessage() {
        return message;
    }

    public void setChallengeConfigs(List<ChallengeConfig> challengeConfigs) {
        this.challengeConfigs = challengeConfigs;
    }

    public List<ChallengeConfig> getChallengeConfigs() {
        return challengeConfigs;
    }
}
