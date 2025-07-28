package release

import "flag"

var (
	dryrun *bool
)

func parseFlags()  {
	dryrun = flag.Bool("dry-run", false, "skip version procedure and release commit/tag")
	flag.Parse()
}
