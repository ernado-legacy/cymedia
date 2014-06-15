package photo

import (
	"io"
)

import (
	"errors"
	"github.com/ernado/cymedia"
	"github.com/ernado/poputchiki-api/weed"
	"github.com/rainycape/magick"
	"log"
	"sync"
	"time"
)

const (
	WEBP        = "webp"
	JPEG        = "jpeg"
	PNG         = "png"
	WEBP_FORMAT = "image/webp"
	JPEG_FORMAT = "image/jpeg"
)

type File struct {
	Fid  string    `json:"fid"           bson:"fid"`
	Time time.Time `json:"time"          bson:"time"`
	Type string    `json:"type"          bson:"type"`
	Size int64     `json:"size"          bson:"size"`
	Url  string    `json:"url,omitempty" bson:"-"`
}

type Photo struct {
	ImageJpeg     File
	ImageWebp     File
	ThumbnailJpeg File
	ThumbnailWebp File
}

type Uploader struct {
	adapter       *weed.Adapter
	maxSize       int
	thumbnailSize int
}

func (uploader *Uploader) upload(image *magick.Image, format string) (string, string, int64, error) {
	encodeReader, encodeWriter := io.Pipe()
	go func() {
		defer encodeWriter.Close()
		info := magick.NewInfo()
		info.SetFormat(format)
		if err := image.Encode(encodeWriter, info); err != nil {
			log.Println(err)
		}
	}()
	return uploader.adapter.Upload(encodeReader, "image", format)
}

func (uploader *Uploader) Upload(length int64, f io.Reader, progress chan float32) (*Photo, error) {
	progressWriter := cymedia.Progress(length, 10, progress)
	im, err := magick.Decode(io.TeeReader(f, progressWriter))
	if err != nil {
		return nil, err
	}

	height, width := float64(im.Height()), float64(im.Width())
	max := float64(uploader.maxSize)
	ratio := max / width
	if height > width {
		ratio = max / height
	}

	if (height < max && width < max) || uploader.maxSize == 0 {
		ratio = 1.0
	}

	// preparing variables for concurrent uploading/processing
	failed := false
	var photoWebp, photoJpeg File
	var purlJpeg, purlWebp string
	var thumbWebp, thumbJpeg File
	var thumbPurlJpeg, thumbPurlWebp string
	wg := new(sync.WaitGroup)
	wg.Add(6)

	// generating abstract upload function
	upload := func(image *magick.Image, url *string, photo *File, extension, format string) {
		defer wg.Done()
		fid, purl, size, err := uploader.upload(image, extension)
		*url = purl
		if err != nil {
			failed = true
			return
		}
		*photo = File{Fid: fid, Time: time.Now(), Type: format, Size: size}
	}

	// resize image and upload to weedfs
	go func() {
		defer wg.Done()
		resized, err := im.Resize(int(width*ratio), int(height*ratio), magick.FBox)
		if err != nil {
			failed = true
			return
		}
		go upload(resized, &purlWebp, &photoWebp, WEBP, WEBP_FORMAT)
		go upload(resized, &purlJpeg, &photoJpeg, JPEG, JPEG_FORMAT)
	}()

	// make thumbnail and upload to weedfs
	go func() {
		defer wg.Done()
		thumbnail, err := im.CropToRatio(1.0, magick.CSCenter)
		if err != nil {
			failed = true
			return
		}
		thumbnail, err = thumbnail.Resize(uploader.thumbnailSize, uploader.thumbnailSize, magick.FBox)
		if err != nil {
			failed = true
			return
		}
		go upload(thumbnail, &thumbPurlWebp, &thumbWebp, WEBP, WEBP_FORMAT)
		go upload(thumbnail, &thumbPurlJpeg, &thumbJpeg, JPEG, JPEG_FORMAT)
	}()
	wg.Wait()

	if failed {
		return nil, errors.New("failed")
	}

	return &Photo{photoJpeg, photoWebp, thumbJpeg, thumbWebp}, nil
}
