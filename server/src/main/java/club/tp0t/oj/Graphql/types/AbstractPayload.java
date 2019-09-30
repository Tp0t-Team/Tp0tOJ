package club.tp0t.oj.Graphql.types;

public abstract class AbstractPayload {
    private String message;

    public AbstractPayload(String message) {
        this.message = message;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }
}
