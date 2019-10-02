package club.tp0t.oj.Service;

import club.tp0t.oj.Dao.ChallengeRepository;
import club.tp0t.oj.Entity.Challenge;
import club.tp0t.oj.Graphql.types.ChallengeInfo;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class ChallengeService {
    @Autowired
    private ChallengeRepository challengeRepository;

    public List<Challenge> getAllChallenges() {
        return challengeRepository.getEnabledChallenges();
    }

}
