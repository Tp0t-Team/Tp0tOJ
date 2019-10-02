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
    private String stuNumber;

    @NotEmpty
    @Column(name = "papssword")
    private String password;

    @NotEmpty
    @Column(name = "department")
    private String department;

    @NotEmpty
    @Column(name = "state")
    private String state;  // normal, disabled, protected

    @NotEmpty
    @Column(name = "QQ")
    private String QQ;

    @NotEmpty
    @Column(name = "mail")
    private String mail;

    @NotNull
    @Column(name = "join_time")
    private Timestamp joinTime;

    @NotEmpty
    @Column(name = "role")
    private String role;
    // admin: administrator
    // member: common user
    // team: user from Tp0t team

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

    @NotEmpty
    @Column(name = "grade")
    private String grade;

    public String getGrade() {
        return grade;
    }

    public long getUserId() {
        return userId;
    }

    public String getName() {
        return name;
    }

    public String getStuNumber() {
        return stuNumber;
    }

    public String getPassword() {
        return password;
    }

    public String getDepartment() {
        return department;
    }

    public String getState() {
        return state;
    }

    public String getQQ() {
        return QQ;
    }

    public String getMail() {
        return mail;
    }

    public Timestamp getJoinTime() {
        return joinTime;
    }

    public String getRole() {
        return role;
    }

    public long getScore() {
        return score;
    }

    public Timestamp getGmtCreated() {
        return gmtCreated;
    }

    public Timestamp getGmtModified() {
        return gmtModified;
    }

    public int getTopRank() {
        return topRank;
    }

    public Timestamp getProtectedTime() {
        return protectedTime;
    }

    public void setUserId(long userId) {
        this.userId = userId;
    }

    public void setName(String name) {
        this.name = name;
    }

    public void setStuNumber(String stuNumber) {
        this.stuNumber = stuNumber;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public void setDepartment(String department) {
        this.department = department;
    }

    public void setState(String state) {
        this.state = state;
    }

    public void setQQ(String QQ) {
        this.QQ = QQ;
    }

    public void setMail(String mail) {
        this.mail = mail;
    }

    public void setJoinTime(Timestamp joinTime) {
        this.joinTime = joinTime;
    }

    public void setRole(String role) {
        this.role = role;
    }

    public void setScore(long score) {
        this.score = score;
    }

    public void setGmtCreated(Timestamp gmtCreated) {
        this.gmtCreated = gmtCreated;
    }

    public void setGmtModified(Timestamp gmtModified) {
        this.gmtModified = gmtModified;
    }

    public void setTopRank(int topRank) {
        this.topRank = topRank;
    }

    public void setProtectedTime(Timestamp protectedTime) {
        this.protectedTime = protectedTime;
    }

    public void setGrade(String grade) {
        this.grade = grade;
    }

}
