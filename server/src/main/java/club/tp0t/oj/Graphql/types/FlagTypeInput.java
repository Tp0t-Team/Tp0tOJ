package club.tp0t.oj.Graphql.types;

public class FlagTypeInput {
    private Boolean dynamic;
    private String value;

    public Boolean getDynamic() {
        return dynamic;
    }

    public String getValue() {
        return value;
    }

    public void setDynamic(Boolean dynamic) {
        this.dynamic = dynamic;
    }

    public void setValue(String value) {
        this.value = value;
    }

    public boolean checkPass() {
        value = value.trim();
        return !value.equals("");
    }
}

