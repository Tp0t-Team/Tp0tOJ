package club.tp0t.oj.Graphql.types;

public class RankResultDesc {
    private String userId;
    private String name;
    private String score;

    public void setUserId(String userId) {
        this.userId = userId;
    }

    public void setName(String name) {
        this.name = name;
    }

    public void setScore(String score) {
        this.score = score;
    }

    public String getUserId() {
        return userId;
    }

    public String getName() {
        return name;
    }

    public String getScore() {
        return score;
    }
}
