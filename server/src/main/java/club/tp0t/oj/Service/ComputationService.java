package club.tp0t.oj.Service;

import club.tp0t.oj.Util.OjConfig;
import org.jetbrains.annotations.NotNull;
import org.springframework.stereotype.Service;
import java.text.SimpleDateFormat;

import java.util.Date;
@Service
public class ComputationService {
    private final Boolean competitionMode;
    private final Boolean registerAllow;
    private final Date beginTime;
    private final Date endTime;
    private final OjConfig ojConfig;
    public ComputationService(@NotNull OjConfig ojConfig){
        this.ojConfig = ojConfig;
        this.competitionMode = ojConfig.isCompetition();
        this.registerAllow = ojConfig.isAllowRegister();
        this.beginTime = ojConfig.getBeginTime();
        this.endTime = ojConfig.getEndTime();
    }

    public String getEndTime() {
        SimpleDateFormat UTC = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSSXXX");
        SimpleDateFormat format = new SimpleDateFormat("yyyy/MM/dd");
        try {
            return UTC.format(endTime);
        }catch (Exception e){
            e.printStackTrace();
        }
        return null;
    }

    public String getBeginTime() {
        SimpleDateFormat UTC = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSSXXX");
        SimpleDateFormat format = new SimpleDateFormat("yyyy/MM/dd");
        try {
            return UTC.format(beginTime);
        }catch (Exception e){
            e.printStackTrace();
        }
        return null;
    }

    public Boolean getCompetitionMode() {
        return competitionMode;
    }

    public Boolean getRegisterAllow() {
        return registerAllow;
    }
    public void setCompetitionMode(boolean competitionMode){
        this.ojConfig.setCompetitionMode(competitionMode);
    }
    public void setRegisterAllow(boolean registerAllow){
        this.ojConfig.setAllowRegister(registerAllow);
    }
    public void setEndTime(String endTime){
        SimpleDateFormat UTC = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSSXXX");
        SimpleDateFormat format = new SimpleDateFormat("yyyy/MM/dd");
        try {
            this.ojConfig.setEndTime(format.parse(endTime));
        }catch (Exception e){
            e.printStackTrace();
        }

    }
    public void setBeginTime(String beginTime){
        SimpleDateFormat UTC = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSSXXX");
        SimpleDateFormat format = new SimpleDateFormat("yyyy/MM/dd");
        try {
            this.ojConfig.setEndTime(format.parse(beginTime));
        }catch (Exception e){
            e.printStackTrace();
        }

    }


}
