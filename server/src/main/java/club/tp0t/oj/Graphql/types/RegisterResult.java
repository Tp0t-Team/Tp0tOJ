package club.tp0t.oj.Graphql.types;

public class RegisterResult {
    private String message;

    public RegisterResult(String message) {
        this.message = message;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }
}
