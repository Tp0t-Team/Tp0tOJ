package club.tp0t.oj.Util;

import java.util.regex.Pattern;

public class CheckHelper {

    private static final String MAIL_PATTERN_STR = "^[_A-Za-z0-9-+]+(.[_A-Za-z0-9-]+)*@" +
            "[A-Za-z0-9-]+(.[A-Za-z0-9]+)*(.[A-Za-z]{2,})$";
    public static final Pattern MAIL_PATTERN = Pattern.compile(MAIL_PATTERN_STR);

    public static boolean checkGrade(String grade) {
        try {
            int gradeInt = Integer.parseInt(grade);
            return (gradeInt >= 2013 && gradeInt <= 9999);
        } catch (NumberFormatException e) {
            return false;
        }
    }

}
