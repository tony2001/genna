package validate

import (
	"strings"

	"github.com/dizzyfool/genna/generators/base"
	"github.com/dizzyfool/genna/util"
)

// Options for generator
type Options struct {
	base.Options

	// Do not replace primary key name to ID
	KeepPK bool
}

// Def fills default values of an options
func (o *Options) Def() {
	o.Options.Def()

	if strings.Trim(o.Package, " ") == "" {
		o.Package = util.DefaultPackage
	}
}
