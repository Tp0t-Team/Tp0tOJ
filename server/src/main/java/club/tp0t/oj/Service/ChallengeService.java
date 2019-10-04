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

    public ChallengeConfiguration getConfiguration(Challenge challenge) {
        String configuration = challenge.getConfiguration();
        ChallengeConfiguration challengeConfiguration = JSON.parseObject(configuration, ChallengeConfiguration.class);
        return challengeConfiguration;
    }

    public Boolean checkIdExistence(String id) {
        Challenge challenge = challengeRepository.findByChallengeId(id);
        return challenge != null;
    }

    public Boolean updateChallenge(ChallengeMutateInput challengeMutate) {
        Challenge challenge = challengeRepository.findByChallengeId(challengeMutate.getChallengeId());

        String configuration = challenge.getConfiguration();
        ChallengeConfiguration challengeConfiguration = JSON.parseObject(configuration, ChallengeConfiguration.class);

        if (!challengeConfiguration.getType().equals(challengeMutate.getType())) return false;

        challengeConfiguration.setName(challengeMutate.getName());
//        if(challengeMutate.getType() != null) challengeConfiguration.setType(challengeMutate.getType());
        challengeConfiguration.setScoreEx(challengeMutate.getScore());
        challengeConfiguration.setFlagEx(challengeMutate.getFlag());
        challengeConfiguration.setDescription(challengeMutate.getDescription());
        challengeConfiguration.setExternalLink(challengeMutate.getExternal_link());
        challengeConfiguration.setHint(challengeMutate.getHint());

        String configurationUpdated = JSON.toJSONString(challengeConfiguration);
        challenge.setConfiguration(configurationUpdated);
        if (challengeMutate.getState() != null) challenge.setState(challengeMutate.getState());

        challengeRepository.save(challenge);
        return true;
    }

    public void updateChallengeBlood(Challenge challenge) {
        challengeRepository.save(challenge);
    }

//    public Boolean checkFormat(ChallengeMutateInput challengemutate){
//
//        String name = challengemutate.getName();
//        String type = challengemutate.getType();
//        ScoreTypeInput score = challengemutate.getScore();
//        FlagTypeInput flag = challengemutate.getFlag();
//        String description = challengemutate.getDescription();
////        List<String> links = challengemutate.getExternal_link();
////        List<String> hints = challengemutate.getHint();
//        if(name==null || type==null || score==null || flag==null || description==null) {
//            return false;
//        }
//
//        name = name.replaceAll("\\s", "");
//        type = type.replaceAll("\\s", "");
//        description = description.replaceAll("\\s", "");
//
//        return true;
//    }

    public Boolean addChallenge(ChallengeMutateInput challengemutate) {

        ChallengeConfiguration challengeConfiguration = new ChallengeConfiguration();
        challengeConfiguration.setName(challengemutate.getName());
        challengeConfiguration.setType(challengemutate.getType());
        challengeConfiguration.setDescription(challengemutate.getDescription());
        challengeConfiguration.setFlagEx(challengemutate.getFlag());
        challengeConfiguration.setScoreEx(challengemutate.getScore());
        challengeConfiguration.setExternalLink(challengemutate.getExternal_link());
        challengeConfiguration.setHint(challengemutate.getHint());

        String configuration = JSON.toJSONString(challengeConfiguration);
        Challenge challenge = new Challenge();
        challenge.setConfiguration(configuration);

        challenge.setState(challengemutate.getState());

        challenge.setGmtCreated(new Timestamp(System.currentTimeMillis()));
        challenge.setGmtModified(new Timestamp(System.currentTimeMillis()));
        challengeRepository.save(challenge);
        return true;
    }

//    public Boolean removeById(String id) {
//        Challenge challenge = challengeRepository.findByChallengeId(id);
//        if (challenge == null) return false;
//        challengeRepository.delete(challenge);
//        return true;
//    }
}
