package club.tp0t.oj.Service;

import club.tp0t.oj.Dao.FlagProxyRepository;
import club.tp0t.oj.Entity.FlagProxy;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;

@Service
public class FlagProxyService {
    private final FlagProxyRepository flagProxyRepository;

    public FlagProxyService(FlagProxyRepository flagProxyRepository) {
        this.flagProxyRepository = flagProxyRepository;
    }

    public List<Long> getPortsByChallengeId(Long challengeId) {
        List<FlagProxy> flagProxyList = flagProxyRepository.findAllByChallengeId(challengeId);
        List<Long> portList = new ArrayList<>();
        for (FlagProxy flagProxy : flagProxyList) {
            portList.add(flagProxy.getPort());
        }
        return portList;
    }

    public String getFlagByChallengeIdAndPort(Long challengeId, Long port) {
        FlagProxy tmpFlagProxy = flagProxyRepository.findByChallengeIdAndPort(challengeId, port);
        if (tmpFlagProxy != null) {
            return tmpFlagProxy.getFlag();
        } else {
            return "flag{no_flag_found_please_contact_organizer}";
        }
    }
}
