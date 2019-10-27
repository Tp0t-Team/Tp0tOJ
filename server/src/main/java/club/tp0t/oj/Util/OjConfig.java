package club.tp0t.oj.Util;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.PropertySource;
import org.springframework.stereotype.Component;

@Component
@PropertySource(value = {"classpath:/application.properties", "file:./application.properties"})
@ConfigurationProperties(prefix = "oj")
public class OjConfig {
    private String host;
    private String name;
    private String salt;

    public String getHost() {
        return host;
    }

    public String getName() {
        return name;
    }

    public String getSalt() {
        return salt;
    }

    public void setHost(String host) {
        this.host = host;
    }

    public void setName(String name) {
        this.name = name;
    }

    public void setSalt(String salt) {
        this.salt = salt;
    }
}
