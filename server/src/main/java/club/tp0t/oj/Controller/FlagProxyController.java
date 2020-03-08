package club.tp0t.oj.Controller;

import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Service.ChallengeService;
import club.tp0t.oj.Service.FlagProxyService;
import club.tp0t.oj.Util.ChallengeConfiguration;
import com.alibaba.fastjson.JSONObject;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.*;

import java.util.ArrayList;
import java.util.List;

@Controller
public class FlagProxyController {
    private final FlagProxyService flagProxyService;
    private final ChallengeService challengeService;

    public FlagProxyController(FlagProxyService flagProxyService,
                               ChallengeService challengeService) {
        this.flagProxyService = flagProxyService;
        this.challengeService = challengeService;
    }

    @RequestMapping(value = "/test", method = RequestMethod.GET)
    @ResponseBody
    public String test() {
        return "flagProxy";
    }

    @RequestMapping(value = "/ports/{challengeId}/{key}",
            method = RequestMethod.GET, produces = "application/json;charset=UTF-8")
    public String proxyPorts(@PathVariable Long challengeId,
                             @PathVariable String key) {
        Challenge tmpChallenge = challengeService.getChallengeByChallengeId(challengeId);
        if (tmpChallenge != null) {  // right challengeId
            ChallengeConfiguration challengeConfig = ChallengeConfiguration.parseConfiguration(tmpChallenge.getConfiguration());
            if (key.equals(challengeConfig.getFlag().getValue()) &&
                    challengeConfig.getFlag().isDynamic()) {  // right key & dynamic flag
                // get ports
                List<Long> portsList = flagProxyService.getPortsByChallengeId(challengeId);
                JSONObject result = new JSONObject();
                result.put("msg", "success");
                result.put("ports", portsList);
                return result.toString();
            } else {  // wrong init flag || flag not dynamic
                JSONObject result = new JSONObject();
                result.put("msg", "auth error");
                result.put("ports", new ArrayList<>());
                return result.toString();
            }
        } else {  // wrong challenge id
            JSONObject result = new JSONObject();
            result.put("msg", "auth error");
            result.put("ports", new ArrayList<>());
            return result.toString();
        }
    }

    @RequestMapping(value = "/flagByPort/{challengeId}/{key}/{port}",
            method = RequestMethod.GET, produces = "application/json;charset=UTF-8")
    public String flagByPort(@PathVariable Long challengeId,
                             @PathVariable String key,
                             @PathVariable Long port) {
        Challenge tmpChallenge = challengeService.getChallengeByChallengeId(challengeId);
        if (tmpChallenge != null) {  // right challengeId
            ChallengeConfiguration challengeConfig = ChallengeConfiguration.parseConfiguration(tmpChallenge.getConfiguration());
            if (key.equals(challengeConfig.getFlag().getValue()) &&
                    challengeConfig.getFlag().isDynamic()) {  // right key & dynamic flag
                String realFlag = flagProxyService.getFlagByChallengeIdAndPort(challengeId, port);
                JSONObject result = new JSONObject();
                result.put("msg", "success");
                result.put("flag", realFlag);
                return result.toString();
            } else {  // wrong key || flag not dynamic
                JSONObject result = new JSONObject();
                result.put("msg", "auth error");
                result.put("flag", "");
                return result.toString();
            }
        } else {  // wrong challengeId
            JSONObject result = new JSONObject();
            result.put("msg", "auth error");
            result.put("flag", "");
            return result.toString();
        }
    }
}
