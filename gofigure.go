// GoFigure is a small utility library for reading configuration files. It's usefuly especially if you want
// to load many files recursively (think /etc/apache2/mods-enabled/*.conf).
//
// It can support multiple formats, as long as you take a file and unmarshal it into a struct containing
// your configurations. Right now the only implemented formats are YAML and JSON files, but feel free to
// add more :)
package gofigure

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Decoder is the interface for config decoders (right now we've just implemented a YAML one)
type Decoder interface {

	// Decode reads data from the io stream, and unmarshals it to the config struct.
	Decode(r io.Reader, config interface{}) error

	// CanDecode should return true if a file is decode-able by the decoder,
	// based on extension or similar mechanisms
	CanDecode(path string) bool
}

// Loader traverses directories recursively and lets the decoder decode relevant files.
//
// It can also explicitly decode single files
type Loader struct {
	decoder Decoder

	// StrictMode determines whether the loader will completely fail on any IO or decoding error,
	// or whether it will continue traversing all files even if one of them is invalid.
	StrictMode bool
}

// NewLoader creates and returns a new Loader wrapping a decoder, using strict mode if specified
func NewLoader(d Decoder, strict bool) *Loader {
	return &Loader{
		decoder:    d,
		StrictMode: strict,
	}
}

// LoadRecursive takes a pointer to a struct containing configurations, and a series of paths.
// It then traverses the paths recursively in their respective order, and lets the decoder decode
// every relevant file.
func (l Loader) LoadRecursive(config interface{}, paths ...string) error {

	ch, cancelc := walk(paths...)
	defer close(cancelc)

	for path := range ch {
		log.Println(path)
		if l.decoder.CanDecode(path) {

			err := l.LoadFile(config, path)
			if err != nil {
				log.Printf("Error loading %s: %s", path, err)
				if l.StrictMode {
					return err
				}
			}

		}
	}

	return nil
}

// LoadFile takes a pointer to a struct containing configurations, and a path to a file,
// and uses the decoder to read the file's contents into the struct. It returns an
// error if the file could not be opened or properly decoded
func (l Loader) LoadFile(config interface{}, path string) error {
	fp, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file %s: %s", path, err)
		if l.StrictMode {
			return err
		}
	}
	defer fp.Close()

	err = l.decoder.Decode(fp, config)
	if err != nil {
		log.Printf("Error decodeing file %s: %s", path, err)
		if l.StrictMode {
			return err
		}
	}
	return nil
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

// walkDir recursively traverses a directory, sending every found file's path to the channel ch.
// If no one is reading from ch, it times out after a second of waiting, and quits
func walkDir(path string, ch chan string, cancelc <-chan struct{}) {

	files, err := ioutil.ReadDir(path)

	if err != nil {
		log.Printf("Could not read path %s: %s", path, err)
		return
	}

	for _, file := range files {
		fullpath := filepath.Join(path, file.Name())
		if isdir, err := isDirectory(fullpath); err != nil {
			log.Printf("Could not stat dir %s: %s", fullpath, err)
			continue
		} else if isdir {
			walkDir(fullpath, ch, cancelc)
			continue
		}

		select {
		case ch <- fullpath:

		case <-cancelc:
			log.Printf("Read canceled")
			return
		}

	}

}

// walk takes a series of paths, and traverses them recursively by order, sending all found files
// in the returned channel. It then closes the channel
func walk(paths ...string) (pathchan <-chan string, cancelchan chan<- struct{}) {

	// we make the channel buffered so it can be filled while the consumer loads files
	ch := make(chan string, 100)
	cancelc := make(chan struct{})

	go func() {
		defer close(ch)
		for _, path := range paths {
			walkDir(path, ch, cancelc)
		}
	}()

	return ch, cancelc

}
