package club.tp0t.oj.Graphql.types;

public class SubmitResult {
    private String message;

    public SubmitResult(String message) {
        this.message = message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public String getMessage() {
        return message;
    }
}
