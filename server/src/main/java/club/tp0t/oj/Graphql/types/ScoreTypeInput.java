package club.tp0t.oj.Graphql.types;

public class ScoreTypeInput {
    private Boolean dynamic;
    private String base_score;

    public Boolean getDynamic() {
        return dynamic;
    }

    public String getBase_score() {
        return base_score;
    }

    public void setBase_score(String base_score) {
        this.base_score = base_score;
    }

    public void setDynamic(Boolean dynamic) {
        this.dynamic = dynamic;
    }

    public boolean checkPass() {
        try {
            return Integer.parseInt(base_score) >= 0;
        }catch (NumberFormatException e) {
            return false;
        }
    }
}
