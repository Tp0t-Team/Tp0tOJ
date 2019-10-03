package club.tp0t.oj.Service;

import club.tp0t.oj.Dao.ChallengeRepository;
import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Graphql.types.ChallengeMutateInput;
import club.tp0t.oj.Graphql.types.FlagTypeInput;
import club.tp0t.oj.Graphql.types.ScoreTypeInput;
import club.tp0t.oj.Util.ChallengeConfiguration;
import com.alibaba.fastjson.JSON;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.sql.Timestamp;
import java.util.List;

@Service
public class ChallengeService {
    @Autowired
    private ChallengeRepository challengeRepository;

    public List<Challenge> getEnabledChallenges() {
        // return challengeRepository.getEnabledChallenges();
        return challengeRepository.findByState("enabled");
    }

    public List<Challenge> getAllChallenges() {
        return challengeRepository.findAll();
    }

    public Challenge getChallengeByChallengeId(long challengeId) {
        return challengeRepository.getOne(challengeId);
    }

    public ChallengeConfiguration getConfiguration(Challenge challenge ) {
        String configuration = challenge.getConfiguration();
        ChallengeConfiguration challengeConfiguration = JSON.parseObject(configuration, ChallengeConfiguration.class);
        return challengeConfiguration;
    }

    public Boolean checkIdExistence(String id){
        Challenge challenge = challengeRepository.findByChallengeId(id);
        return challenge != null;
    }
    public Boolean updateChallenge(ChallengeMutateInput challengeconfig){
        Challenge challenge = challengeRepository.findByChallengeId(challengeconfig.getChallengeId());

        String configuration = challenge.getConfiguration();
        ChallengeConfiguration challengeConfiguration = JSON.parseObject(configuration, ChallengeConfiguration.class);

        if(challengeconfig.getName() != null) challengeConfiguration.setName(challengeconfig.getName());
        if(challengeconfig.getType() != null) challengeConfiguration.setType(challengeconfig.getType());
        if(challengeconfig.getScore()!= null)challengeConfiguration.setScoreEx(challengeconfig.getScore());
        if(challengeconfig.getFlag()!= null) challengeConfiguration.setFlagEx(challengeconfig.getFlag());
        if(challengeconfig.getDescription()!= null) challengeConfiguration.setDescription(challengeconfig.getDescription());
        if(challengeconfig.getExternal_link()!= null) challengeConfiguration.setExternalLink(challengeconfig.getExternal_link());
        if(challengeconfig.getHint()!= null) challengeConfiguration.setHint(challengeconfig.getHint());

        String configurationUpdated = JSON.toJSONString(challengeConfiguration);
        challenge.setConfiguration(configurationUpdated);

//        if (challengeconfig.getState() != null) challenge.setState(challengeconfig.getState());
        challengeRepository.save(challenge);
        return true;
    }
    public Boolean checkFormat(ChallengeMutateInput challengeconfig){

        String name = challengeconfig.getName();
        String type = challengeconfig.getType();
        ScoreTypeInput score = challengeconfig.getScore();
        FlagTypeInput flag = challengeconfig.getFlag();
        String description = challengeconfig.getDescription();
//        List<String> links = challengeconfig.getExternal_link();
//        List<String> hints = challengeconfig.getHint();

        name = name.replaceAll("\\s", "");
        type = type.replaceAll("\\s", "");
        description = description.replaceAll("\\s", "");

        if(name==null || type==null || score==null || flag==null || description==null) {
                return false;
        }
        return true;
    }

    public Boolean addChallenge(ChallengeMutateInput challengeconfig){

        ChallengeConfiguration challengeConfiguration = new ChallengeConfiguration();
        challengeConfiguration.setName(challengeconfig.getName());
        challengeConfiguration.setType(challengeconfig.getType());
        challengeConfiguration.setDescription(challengeconfig.getDescription());
        challengeConfiguration.setFlagEx(challengeconfig.getFlag());
        challengeConfiguration.setScoreEx(challengeconfig.getScore());
        challengeConfiguration.setExternalLink(challengeconfig.getExternal_link());
        challengeConfiguration.setHint(challengeconfig.getHint());

        String configuration = JSON.toJSONString(challengeConfiguration);
        Challenge challenge = new Challenge();
        challenge.setConfiguration(configuration);

//        challenge.setState(challengeconfig.getState());
        challenge.setState("Disabled");

        challenge.setGmtCreated(new Timestamp(System.currentTimeMillis()));
        challenge.setGmtModified(new Timestamp(System.currentTimeMillis()));
        challengeRepository.save(challenge);
        return true;
    }

    public Boolean removeById(String id ){
        Challenge challenge = challengeRepository.findByChallengeId(id);
        if(challenge == null) return false;
        challengeRepository.delete(challenge);
        return true;
    }
}
