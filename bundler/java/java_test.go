package java

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/metrue/fx/utils"
)

func TestJavaBundler(t *testing.T) {
	t.Run("Scaffold", func(t *testing.T) {
		outputDir, err := ioutil.TempDir("", "fx_java")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(outputDir)

		java := New()
		if err := java.Scaffold(outputDir); err != nil {
			t.Fatal(err)
		}

		diff, _, _, err := utils.Diff(outputDir, "./assets")
		if err != nil {
			t.Fatal(err)
		}

		if diff {
			t.Fatalf("%s is not equal with %s", outputDir, "./assets")
		}
	})

	t.Run("BundleSingleFunc", func(t *testing.T) {
		fd, err := ioutil.TempFile("", "fx_func_*.java")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(fd.Name())

		content := `package fx;

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
		`
		err = ioutil.WriteFile(fd.Name(), []byte(content), 0666)
		if err != nil {
			t.Fatal(err)
		}

		outputDir, err := ioutil.TempDir("", "fx_java")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(outputDir)

		java := New()
		if err := java.Bundle(outputDir, fd.Name()); err != nil {
			t.Fatal(err)
		}

		diff, pre, cur, err := utils.Diff("./assets", outputDir)
		if err != nil {
			t.Fatal(err)
		}

		if !diff {
			t.Fatalf("handle functino should be changed")
		}

		if !reflect.DeepEqual(cur, []byte(content)) {
			t.Fatalf("it should be %s but got %s", content, cur)
		}

		preHandleFunc, err := ioutil.ReadFile("./assets/src/main/java/fx/Fx.java")
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(pre, preHandleFunc) {
			{
			}
			t.Fatalf("it should get %s but got %s", preHandleFunc, pre)
		}
	})

	t.Run("BundleFuncAndDeps", func(t *testing.T) {
		fd, err := ioutil.TempFile("", "fx_func_*.js")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(fd.Name())

		content, err := ioutil.ReadFile("./assets/src/main/java/fx/Fx.java")
		if err != nil {
			t.Fatal(err)
		}
		err = ioutil.WriteFile(fd.Name(), content, 0666)
		if err != nil {
			t.Fatal(err)
		}

		addFunc := `
module.exports = (a, b) => a+b
  `
		addFd, err := ioutil.TempFile("", "fx_add_func_*.js")
		if err != nil {
			t.Fatal(err)
		}

		err = ioutil.WriteFile(addFd.Name(), []byte(addFunc), 0644)
		if err != nil {
			t.Fatal(err)
		}

		outputDir, err := ioutil.TempDir("", "fx_java")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(outputDir)

		java := New()
		if err := java.Bundle(outputDir, fd.Name(), addFd.Name()); err != nil {
			t.Fatal(err)
		}

		diff, pre, cur, err := utils.Diff("./assets", outputDir)
		if err != nil {
			t.Fatal(err)
		}

		if !diff {
			t.Fatalf("handle functino should be changed")
		}

		if !reflect.DeepEqual(cur, []byte(addFunc)) {
			t.Fatalf("it should be %s but got %s", content, cur)
		}
		if pre != nil {
			t.Fatal(pre)
		}
	})
}
