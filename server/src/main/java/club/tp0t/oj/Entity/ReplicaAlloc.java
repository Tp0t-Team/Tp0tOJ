package club.tp0t.oj.Entity;

import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import javax.persistence.*;
import javax.validation.constraints.NotNull;
import java.sql.Timestamp;

@Entity
@EntityListeners(AuditingEntityListener.class)
@Table(name = "replica_alloc")
public class ReplicaAlloc {

    @Id
    @Column(name = "replica_alloc_id")
    @GeneratedValue
    private long replicaAllocId;

    @NotNull
    @JoinColumn(name = "user_id", referencedColumnName = "user_id")
    @OneToOne
    private User user;

    @NotNull
    @JoinColumn(name = "replica_id", referencedColumnName = "replica_id")
    @OneToOne
    private Replica replica;

    @NotNull
    @Column(name = "gmt_created")
    @CreatedDate
    private Timestamp gmtCreated;

    @NotNull
    @Column(name = "gmt_modified")
    @LastModifiedDate
    private Timestamp gmtModified;

    public long getReplicaAllocId() {
        return replicaAllocId;
    }

    public User getUser() {
        return user;
    }

    public Replica getReplica() {
        return replica;
    }

    public Timestamp getGmtCreated() {
        return gmtCreated;
    }

    public Timestamp getGmtModified() {
        return gmtModified;
    }

    public void setReplicaAllocId(long replicaAllocId) {
        this.replicaAllocId = replicaAllocId;
    }

    public void setUser(User user) {
        this.user = user;
    }

    public void setReplica(Replica replica) {
        this.replica = replica;
    }

    public void setGmtCreated(Timestamp gmtCreated) {
        this.gmtCreated = gmtCreated;
    }

    public void setGmtModified(Timestamp gmtModified) {
        this.gmtModified = gmtModified;
    }
}
