package cronfab_test

import (
	"fmt"
	"os"
	"time"

	"github.com/aalpar/cronfab"
)

// crontab event at every 5 second interval
func ExampleEveryFiveSeconds() {
	markers, err := cronfab.SecondContabConfig.ParseCronTab("*/5 * * * * *")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("%v\n", markers)

	// run for 2 intervals
	for i := 0; i < 2; i++ {
		t1, err := cronfab.SecondContabConfig.Next(markers, time.Now())
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		dt := t1.Sub(time.Now())
		time.Sleep(dt)
		fmt.Fprintf(os.Stderr, "time: %v\n", time.Now())
	}

	// Output: 0-59/5 0-59/1 0-23/1 1-31/1 1-12/1 0-6/1
}
