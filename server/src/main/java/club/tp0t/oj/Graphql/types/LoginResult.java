package club.tp0t.oj.Graphql.types;

public class LoginResult {
    private String message;
    private String userId;

    public LoginResult(String message, String userId) {
        this.message = message;
        this.userId = userId;
    }
    public LoginResult(String message) {
        this.message = message;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public String getUserId() {
        return userId;
    }
}
