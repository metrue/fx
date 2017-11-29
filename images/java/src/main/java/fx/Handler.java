package fx;

import io.javalin.Javalin;
import org.json.JSONObject;

public class Handler {
    public int handle(JSONObject input) {
        String a = input.get("a").toString();
        String b = input.get("b").toString();
        return Integer.parseInt(a) + Integer.parseInt(b);
    }
}
