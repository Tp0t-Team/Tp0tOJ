package club.tp0t.oj.Graphql.types;

public class SubmitInput {
    private String challengeId;
    private String flag;

    public String getChallengeId() {
        return challengeId;
    }

    public String getFlag() {
        return flag;
    }

    public void setChallengeId(String challengeId) {
        this.challengeId = challengeId;
    }

    public void setFlag(String flag) {
        this.flag = flag;
    }
}
