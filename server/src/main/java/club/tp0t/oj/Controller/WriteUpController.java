package club.tp0t.oj.Controller;

import club.tp0t.oj.Service.WriteupService;
import club.tp0t.oj.Util.CompetitionHelper;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.multipart.MultipartFile;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpSession;
import java.util.Calendar;
import java.util.Date;

@Controller
public class WriteUpController {
    private final WriteupService writeupService;
    private final CompetitionHelper competitionHelper;

    public WriteUpController(WriteupService writeupService, CompetitionHelper competitionHelper) {
        this.writeupService = writeupService;
        this.competitionHelper = competitionHelper;
    }

    @RequestMapping(value = "/writeup")
    @ResponseBody
    public ResponseEntity writeup(@RequestParam("writeup") MultipartFile file, HttpServletRequest request) {
        if (!competitionHelper.getCompetition()) {
            return new ResponseEntity(HttpStatus.FORBIDDEN);
        }
        Date now = new Date();
        Calendar finish = Calendar.getInstance();
        finish.setTime(competitionHelper.getEndTime());
        finish.add(Calendar.HOUR, 1);
        Date finishTime = finish.getTime();
        if (now.compareTo(competitionHelper.getEndTime()) < 0 || now.compareTo(finishTime) > 0) {
            return new ResponseEntity(HttpStatus.FORBIDDEN);
        }
        HttpSession session = request.getSession();
        if (session.getAttribute("isLogin") == null || !(boolean) session.getAttribute("isLogin")) {
            return new ResponseEntity(HttpStatus.UNAUTHORIZED);
        }
        long userId = (Long) request.getSession().getAttribute("userId");
        return writeupService.upload(file, userId);
    }
}
