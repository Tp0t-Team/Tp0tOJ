package club.tp0t.oj.Graphql.types;

public class ForgetResult {
    private  String message;

    public ForgetResult(String message){
        this.message = message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public String getMessage() {
        return message;
    }
}
