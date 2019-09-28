package club.tp0t.oj.Entity;

import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;

import javax.persistence.*;
import javax.validation.constraints.NotNull;
import java.sql.Timestamp;
import java.util.ArrayList;
import java.util.List;

@Entity
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
}
