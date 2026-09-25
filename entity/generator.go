package entity

// GeneratorArgs holds the input parameters for code generation.
type GeneratorArgs struct {
	PkgName       string
	OutName       string
	DocName       string
	Files         []string
	IgnoreImports bool
}
