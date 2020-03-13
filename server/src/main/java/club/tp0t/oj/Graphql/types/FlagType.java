package club.tp0t.oj.Graphql.types;

public class FlagType {
    private Boolean dynamic;
    private String value;
    private int portFrom;
    private int portTo;

    public Boolean getDynamic() {
        return dynamic;
    }

    public String getValue() {
        return value;
    }

    public int getPortFrom() {
        return portFrom;
    }

    public int getPortTo() {
        return portTo;
    }

    public void setDynamic(Boolean dynamic) {
        this.dynamic = dynamic;
    }

    public void setValue(String value) {
        this.value = value;
    }

    public void setPortFrom(int portFrom) {
        this.portFrom = portFrom;
    }

    public void setPortTo(int portTo) {
        this.portTo = portTo;
    }
}

