package club.tp0t.oj.Graphql.types;

import club.tp0t.oj.Util.CheckHelper;

import java.text.DateFormat;
import java.text.SimpleDateFormat;

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

    public boolean checkPass() {
        name = name.replaceAll("\\s", "");
        role = role.replaceAll("\\s", "");
        department = department.replaceAll("\\s", "");
        grade = grade.replaceAll("\\s", "");
        qq = qq.replaceAll("\\s", "");
        mail = mail.replaceAll("\\s", "");
        state = state.replaceAll("\\s", "");

        try {
            DateFormat inputFormat = new SimpleDateFormat("yyyy-MM-dd");
            DateFormat outputFormat = new SimpleDateFormat("yyyy-MM-dd hh:mm:ss");
            protectedTime = outputFormat.format(inputFormat.parse(protectedTime));
        } catch (Exception e) {
            return false;
        }

        return !name.equals("") &&
                checkRole(role) &&
                !department.equals("") &&
                CheckHelper.checkGrade(grade) &&
                !qq.equals("") &&
                !mail.equals("") &&
                CheckHelper.MAIL_PATTERN.matcher(mail).matches() &&
                checkState(state);
    }

    private static boolean checkRole(String role) {
        return role.matches("^(member|team|admin)$");
    }

    private static boolean checkState(String state) {
        return state.matches("^(normal|protected|disabled)$");
    }
}
