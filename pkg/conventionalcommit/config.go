package conventionalcommit

type Config struct {
	MinorTypes []string
	PatchTypes []string
}

var DefaultConfig = Config{
	MinorTypes: defaultMinorTypes,
	PatchTypes: defaultPatchTypes,
}

var defaultMinorTypes = []string{
	"feat",
	"feature",
}

var defaultPatchTypes = []string{
	"fix",
	"perf",
	"performance",
}
