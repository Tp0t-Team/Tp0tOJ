package club.tp0t.oj.Graphql.types;

import java.util.List;

public class ChallengeMutateInput {
    private String challengeId;
    private String name;
    private String type;
    private ScoreTypeInput score;
    private FlagTypeInput flag;
    private String description;
    private List<String> external_link;
    private List<String> hint;
    private String state;

    public void setChallengeId(String challengeId) {
        this.challengeId = challengeId;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public void setExternal_link(List<String> external_link) {
        this.external_link = external_link;
    }

    public void setFlag(FlagTypeInput flag) {
        this.flag = flag;
    }

    public void setHint(List<String> hint) {
        this.hint = hint;
    }

    public void setName(String name) {
        this.name = name;
    }

    public void setScore(ScoreTypeInput score) {
        this.score = score;
    }

    public void setType(String type) {
        this.type = type;
    }

    public String getDescription() {
        return description;
    }

    public ScoreTypeInput getScore() {
        return score;
    }

    public FlagTypeInput getFlag() {
        return flag;
    }

    public List<String> getExternal_link() {
        return external_link;
    }

    public List<String> getHint() {
        return hint;
    }

    public String getName() {
        return name;
    }

    public String getType() {
        return type;
    }

    public String getChallengeId() {
        return challengeId;
    }

    public void setState(String state) {
        this.state = state;
    }

    public String getState() {
        return state;
    }

    public boolean checkPass() {
//        if (name == null || type == null || score == null || flag == null || description == null || external_link == null || hint == null || state == null) return false;
        name = name.trim();
        description = description.trim();
        return !name.equals("") && IsChallengeType(type) && score.checkPass() && flag.checkPass() && !description.equals("") && IsState(state);
    }

    private static boolean IsChallengeType(String type) {
        return type.matches("^(WEB|RE|PWN|MISC|CRYPTO)$");
    }

    private static boolean IsState(String state) {
        return state.matches("^(enabled|disabled)$");
    }
}
