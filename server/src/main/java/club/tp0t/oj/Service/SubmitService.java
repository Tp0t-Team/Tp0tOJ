package club.tp0t.oj.Service;

import club.tp0t.oj.Dao.SubmitRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class SubmitService {
    @Autowired
    private SubmitRepository submitRepository;
}
