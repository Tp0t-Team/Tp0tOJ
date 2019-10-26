package club.tp0t.oj.Graphql.types;

import club.tp0t.oj.Entity.User;

import java.util.ArrayList;
import java.util.List;

public class AllUserInfoResult {

    private List<UserInfo> allUserInfos;

    private String message;

    public AllUserInfoResult(String message) {
        this.message = message;
    }

    // add all information
    public void addAllUserInfo(List<User> users) {
        this.allUserInfos = new ArrayList<>();
        for (int i = 0; i < users.size(); i++) {
            User user = users.get(i);
            this.allUserInfos.add(new UserInfo(user.getName(),
                    user.makeAvatarUrl(),
                    user.getStuNumber(),
                    user.getDepartment(),
                    user.getQq(),
                    user.getMail(),
                    user.getJoinTime(),
                    user.getScore(),
                    user.getTopRank(),
                    user.getProtectedTime(),
                    user.getUserId(),
                    user.getRole(),
                    user.getState(),
                    0,
                    user.getGrade()));
        }
    }

    public List<UserInfo> getAllUserInfos() {
        return allUserInfos;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }


}
