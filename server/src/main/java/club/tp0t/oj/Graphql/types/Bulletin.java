package club.tp0t.oj.Graphql.types;

public class Bulletin {
    private String title;
    private String description;
    private String time;

    public String getDescription() {
        return description;
    }

    public String getTime() {
        return time;
    }

    public String getTitle() {
        return title;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public void setTime(String time) {
        this.time = time;
    }

    public void setTitle(String title) {
        this.title = title;
    }
}
