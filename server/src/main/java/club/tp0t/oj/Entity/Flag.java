package club.tp0t.oj.Entity;

import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;

import javax.persistence.*;
import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.NotNull;
import java.sql.Timestamp;

@Entity
@Table(name = "flag")
public class Flag {

    @Id
    @Column(name = "flag_id")
    @GeneratedValue
    private long flagId;

    @NotNull
    @JoinColumn(name = "challenge_id", referencedColumnName = "challenge_id")
    @ManyToOne
    private Challenge challenge;

    @NotNull
    @JoinColumn(name = "replica_id", referencedColumnName = "replica_id")
    @OneToOne
    private Replica replica;

    @NotEmpty
    @Column(name = "flag")
    private String flag;

    @NotNull
    @Column(name = "gmt_created")
    @CreatedDate
    private Timestamp gmtCreated;

    @NotNull
    @Column(name = "gmt_modified")
    @LastModifiedDate
    private Timestamp gmtModified;

    public long getFlagId() {
        return flagId;
    }

    public Challenge getChallenge() {
        return challenge;
    }

    public Replica getReplica() {
        return replica;
    }

    public String getFlag() {
        return flag;
    }

    public Timestamp getGmtCreated() {
        return gmtCreated;
    }

    public Timestamp getGmtModified() {
        return gmtModified;
    }

    public void setFlagId(long flagId) {
        this.flagId = flagId;
    }

    public void setChallenge(Challenge challenge) {
        this.challenge = challenge;
    }

    public void setReplica(Replica replica) {
        this.replica = replica;
    }

    public void setFlag(String flag) {
        this.flag = flag;
    }

    public void setGmtCreated(Timestamp gmtCreated) {
        this.gmtCreated = gmtCreated;
    }

    public void setGmtModified(Timestamp gmtModified) {
        this.gmtModified = gmtModified;
    }
}
