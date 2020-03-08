package club.tp0t.oj.Graphql.types;

import club.tp0t.oj.Util.OjConfig;

public class CompetitionResult {
    private String message;
    private Boolean competitionMode;
    private Boolean registerAllow;
    private String beginTime;
    private String endTime;
    public CompetitionResult(String message){
        this.message = message;

    }

    public String getMessage() {
        return message;
    }

    public Boolean getCompetitionMode() {
        return competitionMode;
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

    public void setCompetitionMode(Boolean competition) {
        this.competitionMode = competition;
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
