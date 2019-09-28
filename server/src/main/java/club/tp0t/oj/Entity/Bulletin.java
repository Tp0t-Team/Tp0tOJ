package club.tp0t.oj.Entity;

import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;

import javax.persistence.*;
import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.NotNull;
import java.sql.Timestamp;

@Entity
@Table(name = "bulletin")
public class Bulletin {

    @Id
    @Column(name = "bulletin_id")
    @GeneratedValue
    private long bulletinId;

    @NotEmpty
    @Column(name = "content")
    private String content;

    @NotNull
    @Column(name = "topping")
    private boolean topping;

    @NotNull
    @Column(name = "gmt_created")
    @CreatedDate
    private Timestamp gmtCreated;

    @NotNull
    @Column(name = "gmt_modified")
    @LastModifiedDate
    private Timestamp gmtModified;
}
