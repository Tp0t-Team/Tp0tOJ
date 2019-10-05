package club.tp0t.oj.Graphql.types;

public class UserInfoUpdateResult {
    private String message;

    public UserInfoUpdateResult(String message) {
        this.message = message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public String getMessage() {
        return message;
    }
}
