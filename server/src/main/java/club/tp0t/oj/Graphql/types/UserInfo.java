package club.tp0t.oj.Graphql.types;

import java.sql.Timestamp;
import java.text.SimpleDateFormat;

public class UserInfo {
    private String name;
    private String stuNumber;
    private String department;
    private String qq;
    private String mail;
    private String joinTime;
    private String score;
    private int topRank;
    private String protectedTime;
    private String userId;
    private String role;
    private String state;
    private int rank;
    private String grade;

    public UserInfo(String name,
                    String stuNumber,
                    String department,
                    String qq,
                    String mail,
                    Timestamp joinTime,
                    long score,
                    int topRank,
                    Timestamp protectedTime,
                    long userId,
                    String role,
                    String state,
                    int rank,
                    String grade) {
        this.name = name;
        this.stuNumber = stuNumber;
        this.department = department;
        this.qq = qq;
        this.mail = mail;
        this.joinTime = new SimpleDateFormat("yyyy/MM/dd HH:mm:ss").format(joinTime);
        this.score = Long.toString(score);
        this.topRank = topRank;
        this.protectedTime = new SimpleDateFormat("yyyy/MM/dd HH:mm:ss").format(protectedTime);
        this.userId = Long.toString(userId);
        this.role = role;
        this.state = state;
        this.rank = rank;
        this.grade = grade;
    }
    public UserInfo(){};

    public String getName() {
        return name;
    }

    public String getStuNumber() {
        return stuNumber;
    }

    public String getDepartment() {
        return department;
    }

    public String getQq() {
        return qq;
    }

    public String getMail() {
        return mail;
    }

    public String getJoinTime() {
        return joinTime;
    }

    public String getScore() {
        return score;
    }

    public int getTopRank() {
        return topRank;
    }

    public String getUserId() {
        return this.userId;
    }

    public String getProtectedTime() {
        return protectedTime;
    }

    public void setName(String name) {
        this.name = name;
    }

    public void setStuNumber(String stuNumber) {
        this.stuNumber = stuNumber;
    }

    public void setDepartment(String department) {
        this.department = department;
    }

    public void setQq(String qq) {
        this.qq = qq;
    }

    public void setMail(String mail) {
        this.mail = mail;
    }

    public void setJoinTime(String joinTime) {
        this.joinTime = joinTime;
    }

    public void setScore(String score) {
        this.score = score;
    }

    public void setTopRank(int topRank) {
        this.topRank = topRank;
    }

    public void setProtectedTime(String protectedTime) {
        this.protectedTime = protectedTime;
    }

    public void setUserId(String userId) {
        this.userId = userId;
    }

    public void setRole(String role) {
        this.role = role;
    }

    public void setState(String state) {
        this.state = state;
    }

    public void setRank(int rank) {
        this.rank = rank;
    }

    public String getRole() {
        return role;
    }

    public String getState() {
        return state;
    }

    public int getRank() {
        return rank;
    }
}
