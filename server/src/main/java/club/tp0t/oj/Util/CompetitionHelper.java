package club.tp0t.oj.Util;

import org.springframework.stereotype.Component;

import java.text.SimpleDateFormat;
import java.util.Date;

@Component
public class CompetitionHelper {
    private final OjConfig ojConfig;

    public CompetitionHelper(OjConfig ojConfig) {
        this.ojConfig = ojConfig;
    }

    public Date getBeginTime() {
        SimpleDateFormat UTC = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSSXXX");
        try {
            return UTC.parse(ojConfig.getBeginTime());
        } catch (Exception e) {
            e.printStackTrace();
        }
        return null;
    }

    public Date getEndTime() {
        SimpleDateFormat UTC = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSSXXX");
        try {
            return UTC.parse(ojConfig.getEndTime());
        } catch (Exception e) {
            e.printStackTrace();
        }
        return null;
    }

    public Boolean getCompetition() {
        return ojConfig.isCompetition();
    }

    public Boolean getRegisterAllow() {
        return ojConfig.isAllowRegister();
    }

    public void setCompetition(boolean competition) {
        this.ojConfig.setCompetition(competition);
    }

    public void setRegisterAllow(boolean registerAllow) {
        this.ojConfig.setAllowRegister(registerAllow);
    }

    public void setBeginTime(String beginTime) {
        SimpleDateFormat UTC = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSSXXX");
        try {
            UTC.parse(beginTime);
            this.ojConfig.setBeginTime(beginTime);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    public void setEndTime(String endTime) {
        SimpleDateFormat UTC = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSSXXX");
        try {
            UTC.parse(endTime);
            this.ojConfig.setEndTime(endTime);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

}
