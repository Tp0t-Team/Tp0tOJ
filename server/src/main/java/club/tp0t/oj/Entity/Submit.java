package club.tp0t.oj.Entity;

import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;

import javax.persistence.*;
import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.NotNull;
import java.sql.Timestamp;

@Entity
@Table(name = "submit")
public class Submit {

    @Id
    @Column(name = "submit_id")
    @GeneratedValue
    private long submitId;

    @NotNull
    @JoinColumn(name = "user_id", referencedColumnName = "user_id")
    @ManyToOne
    private User user;

    @NotNull
    @Column(name = "mark")
    private int mark;

    @NotNull
    @Column(name = "submit_time")
    private Timestamp submitTime;

    @NotEmpty
    @Column(name = "submit_flag")
    private String submitFlag;

    @NotNull
    @Column(name = "correct")
    private boolean correct;

    @NotNull
    @Column(name = "gmt_created")
    @CreatedDate
    private Timestamp gmtCreated;

    @NotNull
    @Column(name = "gmt_modified")
    @LastModifiedDate
    private Timestamp gmtModified;

    @NotNull
    @JoinColumn(name = "challenge_id", referencedColumnName = "challenge_id")
    @ManyToOne
    private Challenge challenge;

    public long getSubmitId() {
        return submitId;
    }

    public User getUser() {
        return user;
    }

    public int getMark() {
        return mark;
    }

    public Timestamp getSubmitTime() {
        return submitTime;
    }

    public String getSubmitFlag() {
        return submitFlag;
    }

    public boolean isCorrect() {
        return correct;
    }

    public Timestamp getGmtCreated() {
        return gmtCreated;
    }

    public Timestamp getGmtModified() {
        return gmtModified;
    }

    public void setSubmitId(long submitId) {
        this.submitId = submitId;
    }

    public void setUser(User user) {
        this.user = user;
    }

    public void setMark(int mark) {
        this.mark = mark;
    }

    public void setSubmitTime(Timestamp submitTime) {
        this.submitTime = submitTime;
    }

    public void setSubmitFlag(String submitFlag) {
        this.submitFlag = submitFlag;
    }

    public void setCorrect(boolean correct) {
        this.correct = correct;
    }

    public void setGmtCreated(Timestamp gmtCreated) {
        this.gmtCreated = gmtCreated;
    }

    public void setGmtModified(Timestamp gmtModified) {
        this.gmtModified = gmtModified;
    }

    public void setChallenge(Challenge challenge) {
        this.challenge = challenge;
    }

    public Challenge getChallenge() {
        return challenge;
    }
}
