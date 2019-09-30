package club.tp0t.oj.Entity;

import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;

import javax.persistence.*;
import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.NotNull;
import java.sql.Timestamp;

@Entity
@Table(name = "challenge")
public class Challenge {

    @Id
    @Column(name = "challenge_id")
    @GeneratedValue
    private long challengeId;

    @NotEmpty
    @Column(name = "description")
    private String description;

    @Column(name = "first_blood")
    private long firstBlood;

    @Column(name = "second_blood")
    private long secondBlood;

    @Column(name = "third_blood")
    private long thirdBlood;

    @NotEmpty
    @Column(name = "state")
    private String state;

    @NotNull
    @Column(name = "gmt_created")
    @CreatedDate
    private Timestamp gmtCreated;

    @NotNull
    @Column(name = "gmt_modified")
    @LastModifiedDate
    private Timestamp gmtModified;

    public long getChallengeId() {
        return challengeId;
    }

    public String getDescription() {
        return description;
    }

    public long getFirstBlood() {
        return firstBlood;
    }

    public long getSecondBlood() {
        return secondBlood;
    }

    public long getThirdBlood() {
        return thirdBlood;
    }

    public String getState() {
        return state;
    }

    public Timestamp getGmtCreated() {
        return gmtCreated;
    }

    public Timestamp getGmtModified() {
        return gmtModified;
    }

    public void setChallengeId(long challengeId) {
        this.challengeId = challengeId;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public void setFirstBlood(long firstBlood) {
        this.firstBlood = firstBlood;
    }

    public void setSecondBlood(long secondBlood) {
        this.secondBlood = secondBlood;
    }

    public void setThirdBlood(long thirdBlood) {
        this.thirdBlood = thirdBlood;
    }

    public void setState(String state) {
        this.state = state;
    }

    public void setGmtCreated(Timestamp gmtCreated) {
        this.gmtCreated = gmtCreated;
    }

    public void setGmtModified(Timestamp gmtModified) {
        this.gmtModified = gmtModified;
    }
}
