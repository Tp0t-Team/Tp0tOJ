package club.tp0t.oj.Entity;

import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;

import javax.persistence.*;
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
    @Column(name = "gmt_created")
    @CreatedDate
    private Timestamp gmtCreated;

    @NotNull
    @Column(name = "gmt_modified")
    @LastModifiedDate
    private Timestamp gmtModified;

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
}
