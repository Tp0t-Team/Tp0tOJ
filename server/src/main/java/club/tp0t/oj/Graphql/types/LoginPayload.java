package club.tp0t.oj.Graphql.types;

public class LoginPayload {
    private String message;
    private long id;

    public LoginPayload(String message, long id) {
        this.message = message;
        this.id = id;
    }
    public LoginPayload(String message) {
        this.message = message;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }
}
