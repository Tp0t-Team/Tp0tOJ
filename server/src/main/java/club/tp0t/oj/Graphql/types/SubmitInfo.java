package club.tp0t.oj.Graphql.types;

public class SubmitInfo {
    private String submitTime;
    private String challengeName;
    private int mark;

    public void setSubmitTime(String submitTime) {
        this.submitTime = submitTime;
    }

    public void setChallengeName(String challengeName) {
        this.challengeName = challengeName;
    }

    public void setMark(int mark) {
        this.mark = mark;
    }

    public String getSubmitTime() {
        return submitTime;
    }

    public String getChallengeName() {
        return challengeName;
    }

    public int getMark() {
        return mark;
    }
}
