package club.tp0t.oj.Entity;

import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import javax.persistence.*;
import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.NotNull;
import java.sql.Timestamp;

@Entity
@EntityListeners(AuditingEntityListener.class)
@Table(name = "reset_token")
public class ResetToken {

    @Id
    @Column(name = "token_id")
    @GeneratedValue
    private long tokenId;

    @NotEmpty
    @Column(name = "token")
    private String token;

    @NotNull
    @JoinColumn(name = "user_id", referencedColumnName = "user_id")
    @OneToOne
    private User user;

    @NotNull
    @Column(name = "gmt_created")
    @CreatedDate
    private Timestamp gmtCreated;

    @NotNull
    @Column(name = "gmt_modified")
    @LastModifiedDate
    private Timestamp gmtModified;

    public long getTokenId() {
        return tokenId;
    }

    public String getToken() {
        return token;
    }

    public Timestamp getGmtCreated() {
        return gmtCreated;
    }

    public Timestamp getGmtModified() {
        return gmtModified;
    }

    public User getUser() {
        return user;
    }

    public void setTokenId(long tokenId) {
        this.tokenId = tokenId;
    }

    public void setToken(String token) {
        this.token = token;
    }

    public void setGmtCreated(Timestamp gmtCreated) {
        this.gmtCreated = gmtCreated;
    }

    public void setGmtModified(Timestamp gmtModified) {
        this.gmtModified = gmtModified;
    }

    public void setUser(User user) {
        this.user = user;
    }
}
