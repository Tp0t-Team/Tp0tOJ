package club.tp0t.oj.Util;

public class Flag {
    private boolean dynamic;
    private String value;
    private int portFrom;
    private int portTo;

    public boolean isDynamic() {
        return dynamic;
    }

    public String getValue() {
        return value;
    }

    public void setDynamic(boolean dynamic) {
        this.dynamic = dynamic;
    }

    public void setValue(String value) {
        this.value = value;
    }

    public int getPortFrom() {
        return portFrom;
    }

    public int getPortTo() {
        return portTo;
    }

    public void setPortFrom(int portFrom) {
        this.portFrom = portFrom;
    }

    public void setPortTo(int portTo) {
        this.portTo = portTo;
    }
}
