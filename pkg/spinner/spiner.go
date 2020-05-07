package spinner

import (
	"fmt"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/logrusorgru/aurora"
)

var bars map[string]*pb.ProgressBar

func init() {
	bars = make(map[string]*pb.ProgressBar)
}

// Start spinner
func Start(task string) {
	count := 100
	b, ok := bars[task]
	if !ok {
		b = pb.StartNew(count)
		bars[task] = b
	}
	go func() {
		fmt.Printf("Starting %s\n", task)
		for i := 0; i < count; i++ {
			b.Increment()
			time.Sleep(50 * time.Millisecond)
		}
	}()
}

// Stop spinner
func Stop(task string, err error) {
	b, ok := bars[task]
	if ok {
		b.Finish()
	}
	if err != nil {
		fmt.Printf("%s: %s\n", task, aurora.Red("\u2717"))
	}
}
