package club.tp0t.oj.Dao;

import club.tp0t.oj.Entity.Bulletin;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface BulletinRepository extends JpaRepository<Bulletin, Long> {

}
