package service

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
)

type Imagestore interface {
	Save(laptopid string, imagetype string, imagedata bytes.Buffer) (string, error)
}
type imageinfo struct {
	LaptopId string
	Type     string
	Path     string
}
type Imagememory struct {
	mutex       sync.RWMutex
	imagefolder string
	images      map[string]*imageinfo
}

func Newimagestore(imagefolder string) *Imagememory {
	return &Imagememory{
		imagefolder: imagefolder,
		images:      make(map[string]*imageinfo),
	}
}
func (store *Imagememory) Save(laptopid string, imagetype string, imagedata bytes.Buffer) (string, error) {
	imageid, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("cannot generate imageid:%v", err)
	}
	imagepath := fmt.Sprintf("%s/%s%s", store.imagefolder, imageid, imagetype)
	file, err := os.Create(imagepath)
	if err != nil {
		return "", fmt.Errorf("cannnot create imagefile %v", err)
	}
	_, err = imagedata.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("failed to write file")
	}
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.images[imageid.String()] = &imageinfo{
		LaptopId: laptopid,
		Type:     imagetype,
		Path:     imagepath,
	}
	return imageid.String(), nil

}
