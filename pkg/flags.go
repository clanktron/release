package release

import "flag"

var (
	dryrun *bool
	configFile *string
)

func parseFlags()  {
	dryrun = flag.Bool("dry-run", false, "skip version procedure and release commit/tag")
	configFile = flag.String("c", "", "skip version procedure and release commit/tag")
	flag.Parse()
}
