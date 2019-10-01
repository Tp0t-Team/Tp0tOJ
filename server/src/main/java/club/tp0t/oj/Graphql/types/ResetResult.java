package club.tp0t.oj.Graphql.types;

public class ResetResult {
    private String message;

    public ResetResult(String message) {
        this.message = message;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }
}
