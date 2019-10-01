package club.tp0t.oj.Graphql.types;

public class LoginResult {
    private String message;
    private String userId;
    private String role;

    public LoginResult(String message, String userId, String role) {
        this.message = message;
        this.userId = userId;
        this.role = role;
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

    public void setUserId(String userId) {
        this.userId = userId;
    }

    public String getRole() {
        return role;
    }

    public void setRole(String role) {
        this.role = role;
    }
}
