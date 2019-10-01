package club.tp0t.oj.Graphql.types;

public class LoginInput {
    private String stuNumber;
    private String password;

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public void setStuNumber(String stuNumber) {
        this.stuNumber = stuNumber;
    }

    public String getStuNumber() {
        return stuNumber;
    }
}
