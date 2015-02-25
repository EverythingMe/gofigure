package gofigure

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/EverythingMe/gofigure/json"
	"github.com/EverythingMe/gofigure/yaml"
)

type redisConfig struct {
	Server  string `yaml:"server" json:"server"`
	Monitor int    `yaml:"monitor" json:"monitor"`
	Timeout int    `yaml:"timeout" json:"timeout"`
}

type mysqlConfig struct {
	Server   string `yaml:"server" json:"server"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
}

type config struct {
	Redis redisConfig `yaml:"redis"`
	Mysql mysqlConfig `yaml:"mysql"`
}

var expectedConf = config{
	Redis: redisConfig{
		Server:  "localhost:6379",
		Monitor: 1000,
		Timeout: 10,
	},
	Mysql: mysqlConfig{
		Server:   "localhost:3306",
		User:     "root",
		Password: "yeah right :)",
	},
}

func TestYamlLoader(t *testing.T) {

	conf := config{}
	loader := Loader{
		decoder:    yaml.Decoder{},
		StrictMode: true,
	}

	err := loader.LoadRecursive(&conf, "./testdata")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(conf, expectedConf) {
		t.Errorf("Decoded data not as expected: %v", conf)
	}

	err = loader.LoadFile(&conf, "./testdata/test.yaml")
	if err != nil {
		t.Errorf("Error reading single file: %s", err)
	}

}

func TestJsonLoader(t *testing.T) {
	conf := config{}
	loader := Loader{
		decoder:    json.Decoder{},
		StrictMode: true,
	}

	err := loader.LoadRecursive(&conf, "./testdata")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(conf, expectedConf) {
		t.Errorf("Decoded data not as expected: %v", conf)
	}

	err = loader.LoadFile(&conf, "./testdata/test.json")
	if err != nil {
		t.Errorf("Error reading single file: %s", err)
	}

}

func ExampleLoader() {
	// create our configuration container
	var conf = &struct {
		Redis struct {
			Server  string
			Monitor int
			Timeout int
		}
	}{}

	//if we set some default, the loader will override it
	conf.Redis.Server = "localhost:6377"

	// init a loader with a YAML decoder
	loader := Loader{
		decoder:    yaml.Decoder{},
		StrictMode: true,
	}

	// run recursively on the testdata directory
	err := loader.LoadRecursive(conf, "./testdata")
	if err != nil {
		panic(err)
	}

	fmt.Println(conf.Redis.Server)
	//Output: localhost:6379
}
