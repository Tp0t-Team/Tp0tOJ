package club.tp0t.oj.Service;


import club.tp0t.oj.Dao.BulletinRepository;
import club.tp0t.oj.Entity.Bulletin;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class BulletinService {
    @Autowired
    private BulletinRepository bulletinRepository;

    public List<Bulletin> getAllBulletin() { return bulletinRepository.findAll(); }
}
