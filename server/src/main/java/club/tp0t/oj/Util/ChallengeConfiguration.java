package club.tp0t.oj.Util;

import club.tp0t.oj.Graphql.types.FlagTypeInput;
import club.tp0t.oj.Graphql.types.ScoreTypeInput;
import club.tp0t.oj.Graphql.types.FlagType;
import club.tp0t.oj.Graphql.types.ScoreType;

import java.util.List;

public class ChallengeConfiguration {
    private String name;
    private String type;
    private Score score;
    private Flag flag;
    private String description;
    private List<String> externalLink;
    private List<String> hint;


    public String getName() {
        return name;
    }

    public String getType() {
        return type;
    }

//    public int getScore() {
//        return score.getBaseScore();
//    }

    public Score getScore() { return score; }

    public Flag getFlag() {
        return flag;
    }

    public String getDescription() {
        return description;
    }

    public List<String> getExternalLink() {
        return externalLink;
    }

    public List<String> getHint() {
        return hint;
    }

    public void setName(String name) {
        this.name = name;
    }

    public void setType(String type) {
        this.type = type;
    }

    public void setScore(Score score) {
        this.score = score;
    }

    public void setFlag(Flag flag) {
        this.flag = flag;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public void setExternalLink(List<String> externalLink) {
        this.externalLink = externalLink;
    }

    public void setHint(List<String> hint) {
        this.hint = hint;
    }

    public void setFlagEx(FlagTypeInput flagtype){
        this.flag = new Flag();
        flag.setDynamic(flagtype.getDynamic());
        flag.setValue(flagtype.getValue());
    }

    public void setScoreEx(ScoreTypeInput scoretype){
        this.score = new Score();
        score.setBaseScore(Integer.parseInt(scoretype.getBase_score()));
        score.setDynamic(scoretype.getDynamic());
    }

    public ScoreType getScoreEx(){
        ScoreType scoretype =  new ScoreType();
        scoretype.setDynamic(this.getScore().isDynamic());
        scoretype.setBase_score(Integer.toString(this.getScore().getBaseScore()));
        return scoretype;
    }

    public FlagType getFlagEx(){
        FlagType flagtype =  new FlagType();
        flagtype.setDynamic(this.getFlag().isDynamic());
        flagtype.setValue(this.getFlag().getValue());
        return flagtype;
    }
}

class Score {
    private boolean dynamic;
    private int baseScore;

    public boolean isDynamic() {
        return dynamic;
    }

    public int getBaseScore() {
        return baseScore;
    }

    public void setDynamic(boolean dynamic) {
        this.dynamic = dynamic;
    }

    public void setBaseScore(int baseScore) {
        this.baseScore = baseScore;
    }
}
class Flag {
    private boolean dynamic;
    private String value;

    public boolean isDynamic() {
        return dynamic;
    }

    public String getValue() {
        return value;
    }

    public void setDynamic(boolean dynamic) {
        this.dynamic = dynamic;
    }

    public void setValue(String value) {
        this.value = value;
    }
}
