package club.tp0t.oj.Service;

import club.tp0t.oj.Component.ReplicaAllocHelper;
import club.tp0t.oj.Component.ReplicaHelper;
import club.tp0t.oj.Dao.ChallengeRepository;
import club.tp0t.oj.Dao.FlagProxyRepository;
import club.tp0t.oj.Dao.SubmitRepository;
import club.tp0t.oj.Dao.UserRepository;
import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.FlagProxy;
import club.tp0t.oj.Entity.Replica;
import club.tp0t.oj.Entity.User;
import club.tp0t.oj.Graphql.types.BloodInfo;
import club.tp0t.oj.Graphql.types.ChallengeInfo;
import club.tp0t.oj.Graphql.types.ChallengeMutateInput;
import club.tp0t.oj.Util.ChallengeConfiguration;
import club.tp0t.oj.Util.RankHelper;
import com.alibaba.fastjson.JSON;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Isolation;
import org.springframework.transaction.annotation.Transactional;

import java.util.ArrayList;
import java.util.List;

@Service
public class ChallengeService {
    private final ChallengeRepository challengeRepository;
    private final UserRepository userRepository;
    private final SubmitRepository submitRepository;
    private final ReplicaHelper replicaHelper;
    private final ReplicaAllocHelper replicaAllocHelper;
    private final RankHelper rankHelper;
    private final FlagProxyRepository flagProxyRepository;

    public ChallengeService(ChallengeRepository challengeRepository, UserRepository userRepository, SubmitRepository submitRepository, ReplicaHelper replicaHelper, ReplicaAllocHelper replicaAllocHelper, RankHelper rankHelper, FlagProxyRepository flagProxyRepository) {
        this.challengeRepository = challengeRepository;
        this.userRepository = userRepository;
        this.submitRepository = submitRepository;
        this.replicaHelper = replicaHelper;
        this.replicaAllocHelper = replicaAllocHelper;
        this.rankHelper = rankHelper;
        this.flagProxyRepository = flagProxyRepository;
    }

    @Transactional(isolation = Isolation.REPEATABLE_READ) // maybe this level, I'm not sure.
    public List<ChallengeInfo> getChallengeInfoForUser(long userId) {
        User user = userRepository.findByUserId(userId);

        List<Challenge> challengeList = challengeRepository.findByState("enabled");
        ArrayList<ChallengeInfo> result = new ArrayList<>();
        for (Challenge challenge : challengeList) {
            ChallengeInfo challengeInfo = new ChallengeInfo();
            challengeInfo.setChallengeId(Long.toString(challenge.getChallengeId()));

            // set blood
            List<BloodInfo> blood = new ArrayList<>();
            if (challenge.getFirstBlood() != null)
                blood.add(BloodInfo.fromUser(challenge.getFirstBlood()));
            if (challenge.getSecondBlood() != null)
                blood.add(BloodInfo.fromUser(challenge.getSecondBlood()));
            if (challenge.getThirdBlood() != null)
                blood.add(BloodInfo.fromUser(challenge.getThirdBlood()));
            challengeInfo.setBlood(blood);

            // whether solved by user
            challengeInfo.setDone(submitRepository.findByUserAndChallengeAndCorrect(user, challenge, true) != null);

            // parse from description
            ChallengeConfiguration challengeConfiguration = ChallengeConfiguration.parseConfiguration(challenge.getConfiguration());

            // replace flag proxied port
            List<String> externalLinkList = challengeConfiguration.getExternalLink();
            User tmpUser = userRepository.findByUserId(userId);
            if (challengeConfiguration.getFlag().isDynamic() && // only change proxied flag port
                    tmpUser.getRole().equals("member")) {  // do not change for team & admin
                FlagProxy tmpFlagProxy = flagProxyRepository.findByChallengeIdAndUserId(challenge.getChallengeId(), tmpUser.getUserId());
                List<String> replacedExternalLinkList = new ArrayList<>();
                if (tmpFlagProxy != null) {
                    long port = tmpFlagProxy.getPort();
                    for (String tmpExternalLink : externalLinkList) {
                        replacedExternalLinkList.add(tmpExternalLink.replaceAll("端口", Long.toString(port)));
                    }
                } else {  // no record found for this user & challenge
                    for (String tmpExternalLink : externalLinkList) {
                        replacedExternalLinkList.add(tmpExternalLink.replaceAll("端口", "error!"));
                    }
                }
                challengeInfo.setExternalLink(replacedExternalLinkList);
            } else {  // not dynamic or team/admin user
                challengeInfo.setExternalLink(externalLinkList);
            }

            challengeInfo.setDescription(challengeConfiguration.getDescription());
            challengeInfo.setHint(challengeConfiguration.getHint());
            challengeInfo.setType(challengeConfiguration.getType());
            challengeInfo.setName(challengeConfiguration.getName());
            challengeInfo.setScore(Integer.parseInt(challengeConfiguration.getScoreEx().getBase_score()));

            result.add(challengeInfo);
        }
        return result;
    }

    public List<Challenge> getAllChallenges() {
        return challengeRepository.findAll();
    }

    // TODO: this is an utility function.
    public Challenge getChallengeByChallengeId(long challengeId) {
        // TODO: unnecessary
        return challengeRepository.findByChallengeId(challengeId);
    }

    @Transactional(isolation = Isolation.REPEATABLE_READ) // ensure update after read is ok
    public String updateChallenge(ChallengeMutateInput challengeMutate) {
        Challenge challenge = challengeRepository.findByChallengeId(Long.parseLong(challengeMutate.getChallengeId()));

        // ensure challenge exist
        if (challenge == null) return "No such challenge.";

        // unpack JSON data
        String configuration = challenge.getConfiguration();
        ChallengeConfiguration challengeConfiguration = JSON.parseObject(configuration, ChallengeConfiguration.class);

        // ensure type consistency
        if (!challengeConfiguration.getType().equals(challengeMutate.getType())) return "Update Error";

        // update data
        challengeConfiguration.setName(challengeMutate.getName());
        challengeConfiguration.setScoreEx(challengeMutate.getScore());
        challengeConfiguration.setFlagEx(challengeMutate.getFlag());
        challengeConfiguration.setDescription(challengeMutate.getDescription());
        challengeConfiguration.setExternalLink(challengeMutate.getExternal_link());
        challengeConfiguration.setHint(challengeMutate.getHint());

        // pack JSON data
        String configurationUpdated = JSON.toJSONString(challengeConfiguration);

        // update challenge to DB
        challenge.setConfiguration(configurationUpdated);
        if (challengeMutate.getState() != null) challenge.setState(challengeMutate.getState());
        challenge = challengeRepository.save(challenge);

        // update flag for replicas
        List<Replica> replicas = replicaHelper.getReplicaByChallenge(challenge);
        if (replicas != null) {
            for (Replica replica : replicas) {
                // TODO: use saveAll to speed up
                replicaHelper.updateReplicaFlag(replica, challengeConfiguration.getFlag().getValue());
            }
        }

        return "";
    }

    // TODO: this is an utility function.
    public void updateChallengeBlood(Challenge challenge) {
        // TODO: unnecessary
        challengeRepository.save(challenge);
    }

    @Transactional(isolation = Isolation.REPEATABLE_READ) // constraint for alloc replica
    public String addChallenge(ChallengeMutateInput challengeMutate) {
        // pack JSON data
        ChallengeConfiguration challengeConfiguration = new ChallengeConfiguration();
        challengeConfiguration.setName(challengeMutate.getName());
        challengeConfiguration.setType(challengeMutate.getType());
        challengeConfiguration.setDescription(challengeMutate.getDescription());
        challengeConfiguration.setFlagEx(challengeMutate.getFlag());
        challengeConfiguration.setScoreEx(challengeMutate.getScore());
        challengeConfiguration.setExternalLink(challengeMutate.getExternal_link());
        challengeConfiguration.setHint(challengeMutate.getHint());
        String configuration = JSON.toJSONString(challengeConfiguration);

        // add challenge to DB
        Challenge challenge = new Challenge();
        challenge.setConfiguration(configuration);
        challenge.setState(challengeMutate.getState());
        challenge = challengeRepository.save(challenge);

        rankHelper.addChallenge(challenge.getChallengeId(), challengeConfiguration.getScore().getBaseScore());

        // create replicas and allocate to all users
        List<Replica> replicas = replicaHelper.createReplicas(challenge, 1);
        replicaAllocHelper.allocReplicasForAll(replicas);

        return "";
    }

//    public Boolean removeById(String id) {
//        Challenge challenge = challengeRepository.findByChallengeId(id);
//        if (challenge == null) return false;
//        challengeRepository.delete(challenge);
//        return true;
//    }
}
