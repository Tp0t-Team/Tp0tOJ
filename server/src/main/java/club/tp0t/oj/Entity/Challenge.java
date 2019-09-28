package club.tp0t.oj.Entity;

import org.checkerframework.common.aliasing.qual.Unique;
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
}
