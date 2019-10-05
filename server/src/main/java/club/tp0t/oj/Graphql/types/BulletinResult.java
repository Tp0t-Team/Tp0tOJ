package club.tp0t.oj.Graphql.types;

import club.tp0t.oj.Entity.Bulletin;

import java.util.ArrayList;
import java.util.List;

public class BulletinResult {
    private String message;
    private List<Bulletin> bulletins;

    public BulletinResult(String message) {
        this.message = message;
    }

    public List<Bulletin> getBulletin() {
        return bulletins;
    }

    public void setBulletin(List<Bulletin> bulletins) {
        this.bulletins = bulletins;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }
}
