package club.tp0t.oj.Graphql.types;

import club.tp0t.oj.Entity.User;
import club.tp0t.oj.Util.RankHelper;

public class UserInfoResult {

    private UserInfo userInfo;

    private String message;

    public UserInfoResult(String message) {
        this.message = message;
    }


    // hide some attributes
    public void addOthersUserInfo(User user, RankHelper.UserInfo rank) {
        this.userInfo = new UserInfo(user.getName(),
                user.makeAvatarUrl(),
                "",  // stuNumber
                "",  // department
                "",  // qq
                "", // mail
                user.getJoinTime(),
                rank.getScore(),
                user.getTopRank(),
                user.getProtectedTime(),
                user.getUserId(),
                user.getRole(),
                user.getState(),
                rank.getRank(),
                user.getGrade());
    }

    // add all information
    public void addOwnUserInfo(User user, RankHelper.UserInfo rank) {
        this.userInfo = new UserInfo(user.getName(),
                user.makeAvatarUrl(),
                user.getStuNumber(),
                user.getDepartment(),
                user.getQq(),
                user.getMail(),
                user.getJoinTime(),
                rank.getScore(),
                user.getTopRank(),
                user.getProtectedTime(),
                user.getUserId(),
                user.getRole(),
                user.getState(),
                rank.getRank(),
                user.getGrade());
    }

    public UserInfo getUserInfo() {
        return userInfo;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }


}
