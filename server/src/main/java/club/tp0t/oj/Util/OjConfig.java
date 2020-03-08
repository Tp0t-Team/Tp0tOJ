package club.tp0t.oj.Util;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.PropertySource;
import org.springframework.stereotype.Component;

import java.util.Date;
//import java.text.ParseException;
//import java.text.SimpleDateFormat;
//import java.util.Date;

@Component
@PropertySource(value = {"classpath:/application.properties", "file:./application.properties"})
@ConfigurationProperties(prefix = "oj")
public class OjConfig {
    private String host;
    private String name;
    private String salt;
    private double firstBloodPercentage;
    private double secondBloodPercentage;
    private double thirdBloodPercentage;
    private int halfLife;
    private boolean competitionMode;
    private boolean allowRegister;
    private Date beginTime;
    private Date endTime;
    public String getHost() {
        return host;
    }

    public String getName() {
        return name;
    }

    public String getSalt() { return salt; }

    public double getFirstBloodPercentage() {
        return firstBloodPercentage;
    }

    public double getSecondBloodPercentage() {
        return secondBloodPercentage;
    }

    public double getThirdBloodPercentage() {
        return thirdBloodPercentage;
    }

    public int getHalfLife() {
        return halfLife;
    }

    public boolean isCompetitionMode(){ return competitionMode;}

    public boolean isAllowRegister(){ return  allowRegister;}

    public Date getBeginTime() {
        return beginTime;
    }

    public Date getEndTime() {
        return endTime;
    }

    public void setHost(String host) {
        this.host = host;
    }

    public void setName(String name) {
        this.name = name;
    }

    public void setSalt(String salt) {
        this.salt = salt;
    }

    public void setFirstBloodPercentage(double firstBloodPercentage) {
        this.firstBloodPercentage = firstBloodPercentage;
    }

    public void setSecondBloodPercentage(double secondBloodPercentage) {
        this.secondBloodPercentage = secondBloodPercentage;
    }

    public void setThirdBloodPercentage(double thirdBloodPercentage) {
        this.thirdBloodPercentage = thirdBloodPercentage;
    }

    public void setHalfLife(int halfLife) {
        this.halfLife = halfLife;
    }

    public void setCompetitionMode(boolean competitionMode){
        this.competitionMode = allowRegister;
    }

    public void setAllowRegister(boolean allowRegister) {
        this.allowRegister = allowRegister;
    }

    public void setBeginTime(Date beginTime) {
        this.beginTime = beginTime;
    }

    public void setEndTime(Date endTime) {
        this.endTime = endTime;
    }
}

