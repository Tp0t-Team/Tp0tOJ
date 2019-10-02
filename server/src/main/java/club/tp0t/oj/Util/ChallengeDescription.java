package club.tp0t.oj.Util;

import java.util.List;

public class ChallengeDescription {
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

    public int getScore() {
        return score.getBaseScore();
    }

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
