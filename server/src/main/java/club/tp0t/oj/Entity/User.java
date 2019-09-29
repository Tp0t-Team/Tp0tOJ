package club.tp0t.oj.Entity;

import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;

import javax.persistence.*;
import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.NotNull;
import java.sql.Timestamp;

@Entity
@Table(name = "\"user\"")
public class User {

    @Id
    @Column(name = "user_id", length = 20)
    @GeneratedValue
    private long userId;

    @NotEmpty
    @Column(name = "name", length = 100)
    private String name;

    @NotNull
    @NotEmpty
    @Column(name = "stu_number", length = 20)
    private long stuNumber;

    @NotEmpty
    @Column(name = "papssword")
    private String password;

    @NotEmpty
    @Column(name = "department")
    private String department;

    @NotEmpty
    @Column(name = "state")
    private String state;

    @NotEmpty
    @Column(name = "QQ")
    private String QQ;

    @NotEmpty
    @Column(name = "mail")
    private String mail;

    @NotEmpty
    @Column(name = "join_time")
    private Timestamp joinTime;

    @NotEmpty
    @Column(name = "role")
    private String role;

    @NotNull
    @Column(name = "score")
    private long score;

    @NotNull
    @Column(name = "gmt_created")
    @CreatedDate
    private Timestamp gmtCreated;

    @NotNull
    @Column(name = "gmt_modified")
    @LastModifiedDate
    private Timestamp gmtModified;

    @NotNull
    @Column(name = "top_rank")
    private int topRank;

    @NotNull
    @Column(name = "protected_time")
    private Timestamp protectedTime;
}
