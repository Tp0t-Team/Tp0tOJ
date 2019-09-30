package club.tp0t.oj.Graphql.types;

public class RegisterPayload  {
    private String message;

    public RegisterPayload(String message) {
        this.message = message;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }
}
