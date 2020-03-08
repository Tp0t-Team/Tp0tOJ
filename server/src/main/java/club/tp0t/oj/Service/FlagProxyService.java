package club.tp0t.oj.Service;

import club.tp0t.oj.Dao.ChallengeRepository;
import club.tp0t.oj.Dao.FlagProxyRepository;
import club.tp0t.oj.Dao.UserRepository;
import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.FlagProxy;
import club.tp0t.oj.Entity.User;
import club.tp0t.oj.Util.ChallengeConfiguration;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

@Service
public class FlagProxyService {
    private final FlagProxyRepository flagProxyRepository;
    private final ChallengeRepository challengeRepository;
    private final UserRepository userRepository;

    public FlagProxyService(FlagProxyRepository flagProxyRepository, ChallengeRepository challengeRepository, UserRepository userRepository) {
        this.flagProxyRepository = flagProxyRepository;
        this.challengeRepository = challengeRepository;
        this.userRepository = userRepository;
    }

    public List<Long> getPortsByChallengeId(Long challengeId) {
        Challenge tmpChallenge = challengeRepository.findByChallengeId(challengeId);
        List<FlagProxy> flagProxyList = flagProxyRepository.findAllByChallenge(tmpChallenge);
        List<Long> portList = new ArrayList<>();
        for (FlagProxy flagProxy : flagProxyList) {
            portList.add(flagProxy.getPort());
        }
        return portList;
    }

    public String getFlagByChallengeIdAndPort(Long challengeId, Long port) {
        Challenge tmpChallenge = challengeRepository.findByChallengeId(challengeId);
        FlagProxy tmpFlagProxy = flagProxyRepository.findByChallengeAndPort(tmpChallenge, port);
        if (tmpFlagProxy != null) {
            return tmpFlagProxy.getFlag();
        } else {
            return "flag{no_flag_found_please_contact_organizer}";
        }
    }

    // create flags for proxied challenges(state:enabled) & users(role:member)
    // can be used after challenge flag dynamic change from false to true
    // should solve compatibility problems(?)
    public void updateFlagProxies() {
        List<Challenge> challengeList = challengeRepository.findByState("enabled");
        List<Challenge> tmpChallengeList = challengeRepository.findByState("disabled");
        challengeList.addAll(tmpChallengeList);
        List<User> userList = userRepository.findAllByRole("member");
        for (Challenge tmpChallenge : challengeList) {
            ChallengeConfiguration challengeConfiguration = ChallengeConfiguration.parseConfiguration(tmpChallenge.getConfiguration());
            if (challengeConfiguration.getFlag().isDynamic()) {
                for (User tmpUser : userList) {
                    FlagProxy flagProxy = flagProxyRepository.findByChallengeAndUser(tmpChallenge, tmpUser);
                    if (flagProxy == null) {  // create new flagProxy
                        flagProxy = new FlagProxy();
                        flagProxy.setChallenge(tmpChallenge);
                        flagProxy.setUser(tmpUser);
                        flagProxy.setFlag(randomFlag());
                        flagProxy.setPort(randomPort(50000, 65535));  // TODO: add port range to challenge configuration
                        flagProxyRepository.save(flagProxy);
                    }
                }
            }
        }
    }

    // TODO: add proxied flag for new challenge
    public void addChallenge() {

    }

    // TODO: add proxied flag for new user
    public void addUser() {

    }

    // currently generate uuid as random flag
    private String randomFlag() {
        return UUID.randomUUID().toString().replaceAll("-","");
    }

    // get unused port for flag proxy, ranging from 50000 to 65535
    // make sure the challenge server ports are available
    private long randomPort(long start, long end) {
        List<FlagProxy> flagProxyList = flagProxyRepository.findAll();
        List<Long> portList = new ArrayList<>();
        for (FlagProxy tmpFlagProxy : flagProxyList) {
            portList.add(tmpFlagProxy.getPort());
        }
        for (long i = start; i < end; i++) {
            if(!portList.contains(i)) {
                return i;
            }
        }
        return -1; // no enough ports available
    }
}
