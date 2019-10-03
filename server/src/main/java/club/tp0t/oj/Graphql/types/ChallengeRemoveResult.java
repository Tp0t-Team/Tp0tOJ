package club.tp0t.oj.Graphql.types;

public class ChallengeRemoveResult {
    private  String message;

    public ChallengeRemoveResult(String message){
        this.message = message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public String getMessage() {
        return message;
    }
}
