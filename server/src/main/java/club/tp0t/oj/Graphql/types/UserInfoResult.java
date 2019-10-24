package club.tp0t.oj.Graphql.types;

import club.tp0t.oj.Entity.User;
import com.coxautodev.graphql.tools.SchemaParserDictionary;
import org.springframework.context.annotation.Bean;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;

public class UserInfoResult {

    private UserInfo userInfo;

    private String message;

    public UserInfoResult(String message) {
        this.message = message;
    }


    // hide some attributes
    public void addOthersUserInfo(User user, int rank) {
        this.userInfo = new UserInfo(user.getName(),
                "",  // stuNumber
                "",  // department
                "",  // qq
                "", // mail
                user.getJoinTime(),
                user.getScore(),
                user.getTopRank(),
                user.getProtectedTime(),
                user.getUserId(),
                user.getRole(),
                user.getState(),
                rank,
                user.getGrade());
    }

    // add all information
    public void addOwnUserInfo(User user, int rank) {
        this.userInfo = new UserInfo(user.getName(),
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
                rank,
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
