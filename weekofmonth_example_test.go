package cronfab_test

import (
	"fmt"
	"os"
	"time"

	"github.com/aalpar/cronfab"
)

// crontab event at second saturday of every month
func ExampleWeekOfMonth() {
	markers, err := cronfab.SecondContabConfig.ParseCronTab("1 1 1 * 2 * sat")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("%v\n", markers)

	t1 := time.Now()
	// run for 4 intervals
	for i := 0; i < 4; i++ {
		t1, err = cronfab.SecondContabConfig.Next(markers, t1)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stderr, "time: %v\n", t1)
	}

	// Output: 1-1/1 1-1/1 1-1/1 1-31/1 2-2/1 1-12/1 6-6/1
}
