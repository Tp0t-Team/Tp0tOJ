package club.tp0t.oj.Service;

import club.tp0t.oj.Dao.FlagRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class FlagService {
    @Autowired
    private FlagRepository flagRepository;
}
