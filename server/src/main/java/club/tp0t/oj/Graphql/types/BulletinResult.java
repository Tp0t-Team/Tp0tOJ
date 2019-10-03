package club.tp0t.oj.Graphql.types;

import club.tp0t.oj.Entity.Bulletin;

import java.util.List;

public class BulletinResult {
    private String message;
    private List<Bulletin> bulletin;

    public BulletinResult(String message) {
        this.message = message;
    }

    public List<Bulletin> getBulletin() {
        return bulletin;
    }

    public void setBulletin(List<Bulletin> bulletin) {
        this.bulletin = bulletin;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }
}
