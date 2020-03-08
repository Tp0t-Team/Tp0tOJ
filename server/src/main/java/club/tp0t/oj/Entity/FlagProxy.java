package club.tp0t.oj.Entity;

import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import javax.persistence.*;
import javax.validation.constraints.NotNull;

@Entity
@EntityListeners(AuditingEntityListener.class)
@Table(name = "flag_proxy")
public class FlagProxy {
    @Id
    @Column(name = "flag_proxy_id")
    @GeneratedValue
    private long flagProxyId;

    @JoinColumn(name = "challenge_id", referencedColumnName = "challenge_id")
    @ManyToOne
    private Challenge challenge;

    @NotNull
    @JoinColumn(name = "user_id", referencedColumnName = "user_id")
    @ManyToOne
    private User user;

    @NotNull
    @Column(name = "port")
    private long port;

    @NotNull
    @Column(name = "flag")
    private String flag;

    public long getFlagProxyId() {
        return flagProxyId;
    }

    public Challenge getChallenge() {
        return challenge;
    }

    public User getUser() {
        return user;
    }

    public long getPort() {
        return port;
    }

    public String getFlag() {
        return flag;
    }

    public void setFlagProxyId(long flagProxyId) {
        this.flagProxyId = flagProxyId;
    }

    public void setChallenge(Challenge challenge) {
        this.challenge = challenge;
    }

    public void setUser(User user) {
        this.user = user;
    }

    public void setPort(long port) {
        this.port = port;
    }

    public void setFlag(String flag) {
        this.flag = flag;
    }
}
