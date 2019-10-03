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
    @Column(name = "content", columnDefinition = "text")
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

    @NotNull
    @Column(name = "title")
    private String title;

    @NotNull
    @Column(name = "publish_time")
    @CreatedDate
    private Timestamp publishTime;

    public long getBulletinId() {
        return bulletinId;
    }

    public String getContent() {
        return content;
    }

    public boolean isTopping() {
        return topping;
    }

    public Timestamp getGmtCreated() {
        return gmtCreated;
    }

    public Timestamp getGmtModified() {
        return gmtModified;
    }

    public void setBulletinId(long bulletinId) {
        this.bulletinId = bulletinId;
    }

    public void setContent(String content) {
        this.content = content;
    }

    public void setTopping(boolean topping) {
        this.topping = topping;
    }

    public void setGmtCreated(Timestamp gmtCreated) {
        this.gmtCreated = gmtCreated;
    }

    public void setGmtModified(Timestamp gmtModified) {
        this.gmtModified = gmtModified;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public void setPublishTime(Timestamp publishTime) {
        this.publishTime = publishTime;
    }

    public String getTitle() {
        return title;
    }

    public Timestamp getPublishTime() {
        return publishTime;
    }
}
