package fx;

import io.javalin.Javalin;
import org.json.JSONObject;

public class app {
    public static void main(String[] args) {
        Javalin app = Javalin.start(3000);
        Fx handler = new Fx();
        app.post("/", ctx -> {
            JSONObject obj = new JSONObject(ctx.body());
            ctx.result(""+handler.handle(obj));
        });
    }
}
