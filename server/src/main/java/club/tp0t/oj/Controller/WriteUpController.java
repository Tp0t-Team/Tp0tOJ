package club.tp0t.oj.Controller;

import club.tp0t.oj.Service.WriteupService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.multipart.MultipartFile;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpSession;

@Controller
public class WriteUpController {
    private final WriteupService writeupService;

    public WriteUpController(WriteupService writeupService) {
        this.writeupService = writeupService;
    }

    @RequestMapping(value = "/writeup")
    @ResponseBody
    public ResponseEntity writeup(@RequestParam("writeup") MultipartFile file, HttpServletRequest request) {
        HttpSession session = request.getSession();
        if (session.getAttribute("isLogin") == null || !(boolean) session.getAttribute("isLogin")) {
            return new ResponseEntity(HttpStatus.UNAUTHORIZED);
        }
        long userId = (Long) request.getSession().getAttribute("userId");
        return writeupService.upload(file, userId);
    }
}
