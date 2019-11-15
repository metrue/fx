package spinner

import (
	"math/rand"
	"time"

	"github.com/briandowns/spinner"
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
	s.Color(colors[rand.Intn(len(colors))])
	s.Prefix = task + " "
	s.Start()
}

// Stop spinner
func Stop() {
	s.Stop()
}
