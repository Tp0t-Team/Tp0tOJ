package club.tp0t.oj.Service;


import club.tp0t.oj.Dao.BulletinRepository;
import club.tp0t.oj.Entity.Bulletin;
import org.springframework.stereotype.Service;

import java.sql.Timestamp;
import java.util.List;

@Service
public class BulletinService {
    private final BulletinRepository bulletinRepository;

    public BulletinService(BulletinRepository bulletinRepository) {
        this.bulletinRepository = bulletinRepository;
    }

    public List<Bulletin> getAllBulletin() {
        return bulletinRepository.findAll();
    }

    public boolean addBulletin(String title, String content, boolean topping) {
        Bulletin bulletin = new Bulletin();
        bulletin.setTitle(title);
        bulletin.setContent(content);
        bulletin.setPublishTime(new Timestamp(System.currentTimeMillis()));
        bulletin.setTopping(topping);
        bulletinRepository.save(bulletin);
        return true;
    }
}
