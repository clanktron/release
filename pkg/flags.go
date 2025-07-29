package release

import "flag"

var (
	dryrun *bool
	verbose *bool
	configFile *string
	allowUnclean *bool
)

func parseFlags()  {
	dryrun = flag.Bool("dry-run", false, "skip version procedure, commit/tag, and changelog write")
	verbose = flag.Bool("v", false, "verbose")
	configFile = flag.String("c", "", "config file")
	allowUnclean = flag.Bool("a", false, "allow unclean working tree")
	flag.Parse()
}
