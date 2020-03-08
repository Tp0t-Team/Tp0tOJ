package club.tp0t.oj.Service;

import club.tp0t.oj.Util.OjConfig;
import org.springframework.stereotype.Service;

import java.text.SimpleDateFormat;

@Service
public class CompetitionService {
    private final OjConfig ojConfig;

    public CompetitionService(OjConfig ojConfig) {
        this.ojConfig = ojConfig;
    }

    public String getBeginTime() {
        SimpleDateFormat UTC = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSSXXX");
        try {
            return UTC.format(ojConfig.getBeginTime());
        } catch (Exception e) {
            e.printStackTrace();
        }
        return null;
    }

    public String getEndTime() {
        SimpleDateFormat UTC = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSSXXX");
        try {
            return UTC.format(ojConfig.getEndTime());
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
            this.ojConfig.setBeginTime(UTC.parse(beginTime));
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    public void setEndTime(String endTime) {
        SimpleDateFormat UTC = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSSXXX");
        try {
            this.ojConfig.setEndTime(UTC.parse(endTime));
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

}
