package club.tp0t.oj.Graphql.types;

import club.tp0t.oj.Entity.User;

public class BloodInfo {
    private String userId;
    private String name;
    private String avatar;

    public void setUserId(String userId) {
        this.userId = userId;
    }

    public void setName(String name) {
        this.name = name;
    }

    public void setAvatar(String avatar) {
        this.avatar = avatar;
    }

    public String getUserId() {
        return userId;
    }

    public String getName() {
        return name;
    }

    public String getAvatar() {
        return avatar;
    }

    public static BloodInfo fromUser(User user) {
        BloodInfo bloodInfo = new BloodInfo();
        bloodInfo.name = user.getName();
        bloodInfo.avatar = user.makeAvatarUrl();
        bloodInfo.userId = Long.toString(user.getUserId());
        return bloodInfo;
    }
}
