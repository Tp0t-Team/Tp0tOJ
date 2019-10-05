package club.tp0t.oj.Graphql.types;

public class LogoutResult {
    private String message;
    public LogoutResult(String message) {
        this.message = message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public String getMessage() {
        return message;
    }
}
