package club.tp0t.oj.Graphql.types;

public class BulletinPubInput {
    private String title;
    private String content;
    private boolean topping;

    public void setTitle(String title) {
        this.title = title;
    }

    public void setContent(String content) {
        this.content = content;
    }

    public void setTopping(boolean topping) {
        this.topping = topping;
    }

    public String getContent() {
        return content;
    }

    public String getTitle() {
        return title;
    }

    public boolean getTopping() {
        return topping;
    }

    public boolean checkPass() {
        title = title.trim();
        content = content.trim();
        return !title.equals("");
    }
}
