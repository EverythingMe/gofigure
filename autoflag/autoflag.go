// Package autoflag provides automatic command line flags that tell gofigure where to read config files from.
//
// Simply adding an import to:
//    "github.com/EverythingMe/gofigure/autoflag"
// will result in the flags -conf and -confdir being added to your program's flags.
//
// Then you can call autoflag.Load to either load the file in -conf or all files in -confdir.
//
// Note that autoflag.Load will call flag.Parse if you haven't already parsed the flags.
package autoflag

import (
	"errors"
	"flag"

	"github.com/EverythingMe/gofigure"
)

// ConfigDir keeps the value of the -confdir flag if it was set
var ConfigDir string

// ConfigFile keeps the value of the -conf flag if it was set
var ConfigFile string

// init automatically adds the flags to go/flag
func init() {
	flag.StringVar(&ConfigDir, "confdir", "", "If set, recursively read all config files in -confdir")
	flag.StringVar(&ConfigFile, "conf", "", "If set, read just one config file in -conf")
}

// Load either loads the file specified in -conf or the dir in -confdir with loader l to conf
//
// Note that if both are set, we read just the conf file and exit
func Load(l *gofigure.Loader, conf interface{}) error {

	if !flag.Parsed() {
		flag.Parse()
	}

	if ConfigFile != "" {
		return l.LoadFile(conf, ConfigFile)
	}
	if ConfigDir != "" {
		return l.LoadRecursive(conf, ConfigDir)
	}

	return errors.New("gofigure.autoflag: No -conf or -confdir given")
}
