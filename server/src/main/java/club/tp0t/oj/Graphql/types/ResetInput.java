package club.tp0t.oj.Graphql.types;

public class ResetInput {
    private String name;
    private String stuNumber;
    private String mail;

    public String getName() {
        return name;
    }

    public String getStuNumber() {
        return stuNumber;
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

    public void setMail(String mail) {
        this.mail = mail;
    }
}
