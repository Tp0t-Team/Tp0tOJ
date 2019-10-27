package club.tp0t.oj.Graphql.types;

public class ResetInput {
    private String password;
    private String token;

    public String getPassword() {
        return password;
    }

    public String getToken() {
        return token;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public void setToken(String token) {
        this.token = token;
    }

    public boolean checkPass() {
        return !password.equals("") && !token.equals("");
    }
}
