package club.tp0t.oj.Graphql.types;

public class UserInfoUpdateInput {

    private String userId;
    private String name;
    private String role;
    private String department;
    private String grade;
    private String protectedTime;
    private String qq;
    private String mail;
    private String state;

    public void setUserId(String userId) {
        this.userId = userId;
    }

    public void setName(String name) {
        this.name = name;
    }

    public void setRole(String role) {
        this.role = role;
    }

    public void setDepartment(String department) {
        this.department = department;
    }

    public void setGrade(String grade) {
        this.grade = grade;
    }

    public void setProtectedTime(String protectedTime) {
        this.protectedTime = protectedTime;
    }

    public void setQq(String qq) {
        this.qq = qq;
    }

    public void setMail(String mail) {
        this.mail = mail;
    }

    public void setState(String state) {
        this.state = state;
    }

    public String getUserId() {
        return userId;
    }

    public String getName() {
        return name;
    }

    public String getRole() {
        return role;
    }

    public String getDepartment() {
        return department;
    }

    public String getGrade() {
        return grade;
    }

    public String getProtectedTime() {
        return protectedTime;
    }

    public String getQq() {
        return qq;
    }

    public String getMail() {
        return mail;
    }

    public String getState() {
        return state;
    }
}
