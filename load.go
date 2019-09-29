package grue

import (
	"encoding/json"
	"image"
	_ "image/png" // Init png format
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// LoadImages loads JSON images config file and image it refers.
// See theme JSON files for example.
func LoadImages(configFileName string) (ImageSheetConfig, error) {
	res := ImageSheetConfig{}

	cfgFile, err := os.Open(configFileName)
	if err != nil {
		return res, err
	}
	defer cfgFile.Close()

	cfgBytes, err := ioutil.ReadAll(cfgFile)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(cfgBytes, &res)
	if err != nil {
		return res, err
	}

	path := filepath.Dir(configFileName)

	imageFile, err := os.Open(path + "/" + res.File)
	if err != nil {
		return res, err
	}
	defer imageFile.Close()

	res.Atlas, _, err = image.Decode(imageFile)
	if err != nil {
		return res, err
	}

	return res, nil
}

// LoadTTF loads a true type font
func LoadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	face := truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	})

	return face, nil
}
