package club.tp0t.oj.Service;

import club.tp0t.oj.Dao.UserRepository;
import club.tp0t.oj.Entity.User;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import java.io.File;
import java.io.IOException;

@Service
public class WriteupService {
    private final UserRepository userRepository;

    static private final String basePath = "/writeup";

    public WriteupService(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    public ResponseEntity upload(MultipartFile file, long userId) {
        User user = userRepository.findByUserId(userId);
        if(user == null) {
            return new ResponseEntity(HttpStatus.INTERNAL_SERVER_ERROR);
        }
        // save file
        File dest = new File(basePath + user.getStuNumber());
        try {
            file.transferTo(dest);
            return new ResponseEntity(HttpStatus.OK);
        } catch (IllegalStateException | IOException e) {
            e.printStackTrace();
            return new ResponseEntity(HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }
}
