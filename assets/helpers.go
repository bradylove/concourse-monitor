package assets

import (
	"bytes"
	"image"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
)

var (
	CCIconPath string
)

func init() {
	tempPath, err := ioutil.TempDir("", "icons")
	if err != nil {
		log.Fatalf("Failed to create temporary directory for icons: %s", err)
	}

	err = RestoreAsset(tempPath, "icons/cc_icon.png")
	if err != nil {
		log.Fatalf("Failed to restore CC icon: %s", err)
	}
	CCIconPath = filepath.Join(tempPath, "icons", "cc_icon.png")
}

func Reader(name string) io.Reader {
	return bytes.NewReader(MustAsset(name))
}

func Image(name string) image.Image {
	icon, _, err := image.Decode(Reader(name))
	if err != nil {
		log.Fatalf("Failed to decode image: %s", err)
	}

	return icon
}
