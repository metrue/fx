package spinner

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/briandowns/spinner"
	aurora "github.com/logrusorgru/aurora"
)

var s *spinner.Spinner

func init() {
	style := spinner.CharSets[36]
	interval := 100 * time.Millisecond
	s = spinner.New(style, interval)
}

// Start spinner
func Start(task string) {
	colors := []string{
		"red",
		"green",
		"yellow",
		"blue",
		"magenta",
		"cyan",
		"white",
	}

	rand.Seed(time.Now().UnixNano())
	// nolint
	s.Color(colors[rand.Intn(len(colors))])
	s.Prefix = task + " "
	if s.Active() {
		s.Restart()
	} else {
		s.Start()
	}
}

// Stop spinner
func Stop(task string, err error) {
	if err != nil {
		fmt.Println(aurora.Red("\u2717"))
		fmt.Println(aurora.Red("*****************"))
		fmt.Println(err)
		fmt.Println(aurora.Red("*****************"))
	}
	s.Stop()
}
