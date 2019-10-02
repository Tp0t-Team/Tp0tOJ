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
    @Column(name = "description", columnDefinition = "text")
    private String description;

    @JoinColumn(name = "first_blood", referencedColumnName = "user_id")
    @ManyToOne
    private User firstBlood;

    @JoinColumn(name = "second_blood", referencedColumnName = "user_id")
    @ManyToOne
    private User secondBlood;

    @JoinColumn(name = "third_blood", referencedColumnName = "user_id")
    @ManyToOne
    private User thirdBlood;

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

    public void setFirstBlood(User firstBlood) {
        this.firstBlood = firstBlood;
    }

    public void setSecondBlood(User secondBlood) {
        this.secondBlood = secondBlood;
    }

    public User getFirstBlood() {
        return firstBlood;
    }

    public User getSecondBlood() {
        return secondBlood;
    }

    public void setThirdBlood(User thirdBlood) {
        this.thirdBlood = thirdBlood;
    }

    public User getThirdBlood() {
        return thirdBlood;
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
