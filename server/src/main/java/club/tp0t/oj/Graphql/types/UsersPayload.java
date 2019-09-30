package club.tp0t.oj.Graphql.types;

import club.tp0t.oj.Entity.User;
import com.coxautodev.graphql.tools.SchemaParserDictionary;
import org.springframework.context.annotation.Bean;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;

public class UsersPayload {

    private List<UserInfo> userInfos = new ArrayList<>();

    private String message;

    public UsersPayload(String message) {
        this.message = message;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public void addNormalUserInfo(List<User> users) {
        for(int i=0;i<users.size();i++) {
            User tmpUser = users.get(i);
            UserInfo tmpUserInfo = new UserInfo(tmpUser.getName(),
                    "",  // stuNumber
                    "",  // department
                    "",  // qq
                    "", // mail
                    tmpUser.getJoinTime(),
                    tmpUser.getScore(),
                    tmpUser.getTopRank(),
                    tmpUser.getProtectedTime());

        }
    }

    public void addAllUserInfo(List<User> users) {
        for(int i=0;i<users.size();i++) {
            User tmpUser = users.get(i);
            UserInfo tmpUserInfo = new UserInfo(tmpUser.getName(),
                    tmpUser.getStuNumber(),
                    tmpUser.getDepartment(),
                    tmpUser.getQQ(),
                    tmpUser.getMail(),
                    tmpUser.getJoinTime(),
                    tmpUser.getScore(),
                    tmpUser.getTopRank(),
                    tmpUser.getProtectedTime());

        }
    }

    public List<UserInfo> getUserInfos() {
        return userInfos;
    }

}
