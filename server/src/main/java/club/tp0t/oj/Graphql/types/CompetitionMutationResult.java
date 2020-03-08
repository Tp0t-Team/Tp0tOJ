package club.tp0t.oj.Graphql.types;

public class CompetitionMutationResult {
    private String Message;

    public CompetitionMutationResult(String message){
        this.Message = message;
    }
    public String getMessage() {
        return Message;
    }

    public void setMessage(String message) {
        Message = message;
    }
}
