package club.tp0t.oj.Entity;

import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;

import javax.persistence.*;
import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.NotNull;
import java.sql.Timestamp;

@Entity
@Table(name = "replica")
public class Replica {

    @Id
    @GeneratedValue
    @Column(name = "replica_id")
    private long replicaId;

    @NotNull
    @JoinColumn(name = "challenge_id", referencedColumnName = "challenge_id")
    @ManyToOne
    private Challenge challenge;

    @NotEmpty
    @Column(name = "flag", columnDefinition = "text")
    private String flag;

    @NotNull
    @Column(name = "gmt_created")
    @CreatedDate
    private Timestamp gmtCreated;

    @NotNull
    @Column(name = "gmt_modified")
    @LastModifiedDate
    private Timestamp gmtModified;

    @NotNull
    @Column(name = "flag", columnDefinition = "text")
    private String flag;

    public long getReplicaId() {
        return replicaId;
    }

    public Timestamp getGmtCreated() {
        return gmtCreated;
    }

    public Timestamp getGmtModified() {
        return gmtModified;
    }

    public void setReplicaId(long replicaId) {
        this.replicaId = replicaId;
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

    public String getFlag() {
        return flag;
    }

    public void setFlag(String flag) {
        this.flag = flag;
    }

    public Challenge getChallenge() {
        return challenge;
    }

    public void setFlag(String flag) {
        this.flag = flag;
    }

    public String getFlag() {
        return flag;
    }
}
