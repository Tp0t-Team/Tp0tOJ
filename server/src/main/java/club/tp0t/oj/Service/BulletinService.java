package club.tp0t.oj.Service;

import club.tp0t.oj.Dao.BulletinRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class BulletinService {
    @Autowired
    private BulletinRepository bulletinRepository;
}
