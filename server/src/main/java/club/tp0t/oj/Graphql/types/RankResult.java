package club.tp0t.oj.Graphql.types;

import club.tp0t.oj.Entity.User;

import java.util.ArrayList;
import java.util.List;

public class RankResult {
    private String message;
    private List<UserInfo> userInfos = new ArrayList<>();
    public RankResult(String message) {
        this.message = message;
    }

    public void addUserInfos(List<User> users) {
        for(int i=0;i<users.size();i++) {
            User tmpUser = users.get(i);
            UserInfo userInfo = new UserInfo();

            userInfo.setUserId(Long.toString(tmpUser.getUserId()));
            userInfo.setName(tmpUser.getName());
            //userInfo.setStuNumber(tmpUser.getStuNumber());
            userInfo.setStuNumber("");
            userInfo.setDepartment(tmpUser.getDepartment());
            //userInfo.setQq(tmpUser.getQQ());
            userInfo.setQq("");
            //userInfo.setMail(tmpUser.getMail());
            userInfo.setMail("");
            userInfo.setJoinTime(tmpUser.getJoinTime().toString());
            userInfo.setScore(Long.toString(tmpUser.getScore()));
            userInfo.setTopRank(tmpUser.getTopRank());
            userInfo.setProtectedTime(tmpUser.getProtectedTime().toString());

            this.userInfos.add(userInfo);
        }
    }
}
