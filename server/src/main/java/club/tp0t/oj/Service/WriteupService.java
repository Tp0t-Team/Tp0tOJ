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
    private boolean CheckContentType(String contentType){
        boolean flag = false;
        //allowed contentType list
        String[] allowTypeList={"application/pdf","application/msword","application/vnd.openxmlformats-officedocument.wordprocessingml.document","text/markdown","text/plain"};
        for (String allow:allowTypeList) {
            if (allow.compareTo(contentType) == 0) {
                flag = true;
            }
        }
        return flag;
    }
    public ResponseEntity upload(MultipartFile file, long userId) {
        User user = userRepository.findByUserId(userId);
        if(user == null) {
            return new ResponseEntity(HttpStatus.INTERNAL_SERVER_ERROR);
        }
        //check file
        if(file.isEmpty() || !CheckContentType(file.getContentType())){
            return new ResponseEntity(HttpStatus.FORBIDDEN);
        }
        // save file
        String upload_filename = file.getOriginalFilename();
        //destPath is a file name
        String destPath = basePath +"/"+ upload_filename;
        File dest = new File(basePath);
        try {
            file.transferTo(dest);
            return new ResponseEntity(HttpStatus.OK);
        } catch (IllegalStateException | IOException e) {
            e.printStackTrace();
            return new ResponseEntity(HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }
}
