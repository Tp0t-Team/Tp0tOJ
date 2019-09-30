package club.tp0t.oj.Graphql.types;

public class ResetPayload {
    private String message;

    public ResetPayload(String message) {
        this.message = message;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }
}
