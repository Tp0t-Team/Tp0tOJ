package club.tp0t.oj.Graphql.types;

public class CompetitionResult {
    private String message;
    private Boolean competition;
    private Boolean registerAllow;
    private String beginTime;
    private String endTime;

    public String getMessage() {
        return message;
    }

    public Boolean getCompetition() {
        return competition;
    }

    public Boolean getRegisterAllow() {
        return registerAllow;
    }

    public String getBeginTime() {
        return beginTime;
    }

    public String getEndTime() {
        return endTime;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public void setCompetition(Boolean competition) {
        this.competition = competition;
    }

    public void setRegisterAllow(Boolean registerAllow) {
        this.registerAllow = registerAllow;
    }

    public void setBeginTime(String beginTime) {
        this.beginTime = beginTime;
    }

    public void setEndTime(String endTime) {
        this.endTime = endTime;
    }
}
