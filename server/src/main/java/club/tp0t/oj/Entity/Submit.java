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
}
