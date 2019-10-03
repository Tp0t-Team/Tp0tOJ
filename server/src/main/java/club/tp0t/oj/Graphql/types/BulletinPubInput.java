package club.tp0t.oj.Graphql.types;

public class BulletinPubInput {
    private String title;
    private String content;
    private String topping;

    public void setTitle(String title) {
        this.title = title;
    }

    public void setContent(String content) {
        this.content = content;
    }

    public void setTopping(String topping) {
        this.topping = topping;
    }

    public String getContent() {
        return content;
    }

    public String getTitle() {
        return title;
    }

    public String getTopping() {
        return topping;
    }
}
