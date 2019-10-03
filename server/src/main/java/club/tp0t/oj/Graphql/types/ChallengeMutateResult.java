package club.tp0t.oj.Graphql.types;

public class ChallengeMutateResult {
    private String message;

    public ChallengeMutateResult(String message) {
        this.message = message;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }
}
