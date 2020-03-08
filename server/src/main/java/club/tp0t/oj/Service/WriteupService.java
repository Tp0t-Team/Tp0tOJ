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

    private boolean CheckContentType(String contentType) {
        boolean flag = false;
        //allowed contentType list
        String[] allowTypeList = {"application/pdf", "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "text/markdown", "text/plain"};
        for (String allow : allowTypeList) {
            if (allow.compareTo(contentType) == 0) {
                flag = true;
                break;
            }
        }
        return flag;
    }

    public ResponseEntity upload(MultipartFile file, long userId) {
        ResponseEntity result;
        User user = userRepository.findByUserId(userId);
        if (user == null) {
            result = new ResponseEntity(HttpStatus.INTERNAL_SERVER_ERROR);
        } else if (file.isEmpty() || !CheckContentType(file.getContentType())) {
            result = new ResponseEntity(HttpStatus.FORBIDDEN);
        } else {// save file
            String upload_filename = user.getStuNumber() + "-" + user.getName() + "-" + file.getOriginalFilename();//destPath is a file name
            String destPath = basePath + "/" + upload_filename;
            File dest = new File(destPath);
            try {
                file.transferTo(dest);
                result = new ResponseEntity(HttpStatus.OK);
            } catch (IllegalStateException | IOException e) {
                e.printStackTrace();
                result = new ResponseEntity(HttpStatus.INTERNAL_SERVER_ERROR);
            }
        }
        return result;
    }
}
