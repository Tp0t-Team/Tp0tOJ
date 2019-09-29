package club.tp0t.oj.Graphql.types;

public class RegisterInput {
    private String name;
    private String stuNumber;
    private String password;
    private String department;
    private String qq;
    private String mail;

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

    public String getQq() {
        return qq;
    }

    public String getMail() {
        return mail;
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

    public void setQq(String qq) {
        this.qq = qq;
    }

    public void setMail(String mail) {
        this.mail = mail;
    }
}
