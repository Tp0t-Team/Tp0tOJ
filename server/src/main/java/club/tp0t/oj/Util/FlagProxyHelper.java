package club.tp0t.oj.Util;

import club.tp0t.oj.Dao.ChallengeRepository;
import club.tp0t.oj.Dao.FlagProxyRepository;
import club.tp0t.oj.Dao.UserRepository;
import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Entity.FlagProxy;
import club.tp0t.oj.Entity.User;
import club.tp0t.oj.Util.ChallengeConfiguration;
import org.springframework.stereotype.Component;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

@Component
public class FlagProxyHelper {
    private final FlagProxyRepository flagProxyRepository;
    private final ChallengeRepository challengeRepository;
    private final UserRepository userRepository;

    public FlagProxyHelper(FlagProxyRepository flagProxyRepository, ChallengeRepository challengeRepository, UserRepository userRepository) {
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
            return "No flag found";
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
        List<User> tmpUserList = userRepository.findAllByRole("team");
        userList.addAll(tmpUserList);
        tmpUserList = userRepository.findAllByRole("admin");
        userList.addAll(tmpUserList);

        for (Challenge tmpChallenge : challengeList) {
            ChallengeConfiguration challengeConfiguration = ChallengeConfiguration.parseConfiguration(tmpChallenge.getConfiguration());
            if (challengeConfiguration.getFlag().isDynamic()) {  // proxied challenge
                for (User tmpUser : userList) {
                    FlagProxy flagProxy = flagProxyRepository.findByChallengeAndUser(tmpChallenge, tmpUser);
                    if (flagProxy == null) {  // create new flagProxy
                        flagProxy = new FlagProxy();
                        flagProxy.setChallenge(tmpChallenge);
                        flagProxy.setUser(tmpUser);
                        flagProxy.setFlag(randomFlag());
                        flagProxy.setPort(randomPort(challengeConfiguration.getFlag().getPortFrom(), challengeConfiguration.getFlag().getPortTo()));
                        flagProxyRepository.save(flagProxy);
                    }
                }
            }
        }
    }

    // add proxied flag for new challenge
    public void updateChallenge(Challenge challenge) {
        ChallengeConfiguration challengeConfiguration = ChallengeConfiguration.parseConfiguration(challenge.getConfiguration());
        if (challengeConfiguration.getFlag().isDynamic()) {  // to proxied challenge
            List<FlagProxy> flagProxyList = flagProxyRepository.findAllByChallenge(challenge);  // get old records
            if (flagProxyList.size() == 0) {  // not proxied to proxied
                List<User> userList = userRepository.findAllByRole("member");
                List<User> tmpUserList = userRepository.findAllByRole("team");
                userList.addAll(tmpUserList);
                tmpUserList = userRepository.findAllByRole("admin");
                userList.addAll(tmpUserList);

                for (User tmpUser: userList) {
                    FlagProxy flagProxy = new FlagProxy();
                    flagProxy.setChallenge(challenge);
                    flagProxy.setUser(tmpUser);
                    flagProxy.setFlag(randomFlag());
                    flagProxy.setPort(randomPort(challengeConfiguration.getFlag().getPortFrom(), challengeConfiguration.getFlag().getPortTo()));
                    flagProxyRepository.save(flagProxy);
                }

            } else {  // proxied to proxied
                // reallocate ports
                for (FlagProxy flagProxy : flagProxyList) {
                    flagProxy.setPort(randomPort(challengeConfiguration.getFlag().getPortFrom(), challengeConfiguration.getFlag().getPortTo()));
                    flagProxyRepository.save(flagProxy);
                }
            }
        } else {  // to not proxied
            // remove all records
            List<FlagProxy> flagProxies = flagProxyRepository.findAllByChallenge(challenge);
            for (FlagProxy flagProxy : flagProxies) {
                flagProxyRepository.delete(flagProxy);
            }
        }
    }

    // add proxied flag for new user
    public void addUser(User user) {
        List<Challenge> challengeList = challengeRepository.findByState("enabled");
        List<Challenge> tmpChallengeList = challengeRepository.findByState("disabled");
        challengeList.addAll(tmpChallengeList);

        for (Challenge tmpChallenge : challengeList) {
            ChallengeConfiguration challengeConfiguration = ChallengeConfiguration.parseConfiguration(tmpChallenge.getConfiguration());
            if (challengeConfiguration.getFlag().isDynamic()) {
                FlagProxy flagProxy = new FlagProxy();
                flagProxy.setChallenge(tmpChallenge);
                flagProxy.setUser(user);
                flagProxy.setFlag(randomFlag());
                flagProxy.setPort(randomPort(challengeConfiguration.getFlag().getPortFrom(), challengeConfiguration.getFlag().getPortTo()));
                System.out.println("add proxied flag for user :" + challengeConfiguration.getName() + " - " + flagProxy.getFlag());
                flagProxyRepository.save(flagProxy);
            }
        }
    }

    // currently generate uuid as random flag
    private String randomFlag() {
        return UUID.randomUUID().toString().replaceAll("-","");
    }

    // get unused port for flag proxy, ranging from start to end
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

        return (long)(Math.random() * ((start - end) + 1)) + start; // no enough ports available
    }
}