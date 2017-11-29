package fx;

import io.javalin.Javalin;
import org.json.JSONObject;

public class fx {
    public static void main(String[] args) {
        Javalin app = Javalin.start(3000);
        Handler handler = new Handler();
        app.post("/", ctx -> {
            JSONObject obj = new JSONObject(ctx.body());
            ctx.result(""+handler.handle(obj));
        });
    }
}
