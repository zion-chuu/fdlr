package errorHandle

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func ExitWithError(err error) {
	if err != nil {
		// if you are developing, you can use log.Printf("%+v\n", err) to show the trace stack
		fmt.Printf("%v\n", errors.Cause(err))
		os.Exit(1)
	}
}
