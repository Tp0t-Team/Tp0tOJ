package club.tp0t.oj.Graphql.resolvers;

import club.tp0t.oj.Graphql.types.BulletinResult;
import club.tp0t.oj.Service.BulletinService;
import com.coxautodev.graphql.tools.GraphQLQueryResolver;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

@Component
public class BulletinQuery implements GraphQLQueryResolver {
    private final BulletinService bulletinService;

    public BulletinQuery(BulletinService bulletinService) {
        this.bulletinService = bulletinService;
    }

    public BulletinResult allBulletin() {
        BulletinResult bulletinResult = new BulletinResult("");
        bulletinResult.setBulletin(bulletinService.getAllBulletin());
        return bulletinResult;
    }
}
