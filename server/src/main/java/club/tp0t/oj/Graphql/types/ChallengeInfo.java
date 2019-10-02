package club.tp0t.oj.Graphql.types;

import java.util.ArrayList;
import java.util.List;

public class ChallengeInfo {
    private String type;
    private String name;
    private long score;
    private String description;
    private List<String> externalLink = new ArrayList<>();
    private List<String> hint = new ArrayList<>();
    private List<String> blood = new ArrayList<>();
    private boolean done;

    public void setType(String type) {
        this.type = type;
    }

    public void setName(String name) {
        this.name = name;
    }

    public void setScore(long score) {
        this.score = score;
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

    public void setBlood(List<String> blood) {
        this.blood = blood;
    }

    public void setDone(boolean done) {
        this.done = done;
    }

    public String getType() {
        return type;
    }

    public String getName() {
        return name;
    }

    public long getScore() {
        return score;
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

    public List<String> getBlood() {
        return blood;
    }

    public boolean isDone() {
        return done;
    }
}
