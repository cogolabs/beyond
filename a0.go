package beyond

import (
	"log"
	"time"

	"github.com/cogolabs/wait"
)

func init() {
	// prepend file:lineno
	log.SetFlags(log.Flags() | log.Lshortfile)

	// wait for networking
	wait.ForNetwork(5, time.Second)
}
