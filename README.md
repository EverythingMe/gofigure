# GoFigure

GoFigure is a small utility library for reading configuration files. It's usefuly especially if you want to load many files recursively (think `/etc/apache2/mods-enabled/*.conf` on Ubuntu).

It can support multiple formats, as long as you take a file and unmarshal it into a struct containing your configurations. 

Right now the only implemented formats are YAML and JSON files, but feel free to add more :)

## Example usage:

```go 

import (
	"fmt"
	"github.com/EverythingMe/gofigure"
	"github.com/EverythingMe/gofigure/yaml"
)

// this is our configuration container
var conf = &struct {
	Redis struct {
		Server  string
		Monitor int
		Timeout int
	}
}{}

func main() {
	
	//if we set some default, the loader will override it
	conf.Redis.Server = "localhost:6377"

	// init a loader with a YAML decoder
	loader := gofigure.Loader{
		decoder:    yaml.Decoder{},
		StrictMode: true,
	}

	// run recursively on some directory
	err := loader.LoadRecursive(conf, "/etc/myservice/conf.d")
	if err != nil {
		panic(err)
	}

	fmt.Println(conf.Redis.Server)
	//Output: localhost:6379
}

```