package photo

import (
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

type StorageAdapter interface {
	GetUrl(fid string) (url string, err error)
	Upload(reader io.Reader, t, format string) (fid string, purl string, size int64, err error)
}

type Result interface {
	Image() string
	Thumbnail() string
}

type ImageUploader interface {
	Upload(r io.Reader) (Result, error)
}

type UploaderResult struct {
	image     string
	thumbnail string
}

func (u UploaderResult) Image() string {
	return u.image
}

func (u UploaderResult) Thumbnail() string {
	return u.thumbnail
}

func (u *Uploader) EncodeAndUpload(m image.Image) (string, error) {
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		if err := jpeg.Encode(w, m, nil); err != nil {
			w.CloseWithError(err)
		}
	}()
	fid, _, _, err := u.adapter.Upload(r, "image", "jpeg")
	return fid, err
}

func (u *Uploader) Upload(r io.Reader) (Result, error) {
	m, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	// getting image size
	size := m.Bounds().Size()
	width, height := uint(size.X), uint(size.Y)

	// resizing image
	resized := resize.Thumbnail(u.maxSize, u.maxSize, m, resize.Lanczos3)

	// making thumbnail
	if width > height {
		width = uint(float64(u.thumbnailSize*width) / float64(height))
		height = u.thumbnailSize
	} else {
		height = uint(float64(u.thumbnailSize*height) / float64(width))
		width = u.thumbnailSize
	}
	thumbnail := resize.Resize(width, height, m, resize.Lanczos3)
	thumbnail, err = cutter.Crop(thumbnail, cutter.Config{
		Width:  int(u.thumbnailSize),
		Height: int(u.thumbnailSize),
		Mode:   cutter.Centered,
	})
	if err != nil {
		return nil, err
	}

	// uploading
	resizedFid, err := u.EncodeAndUpload(resized)
	if err != nil {
		return nil, err
	}
	thumbnailFid, err := u.EncodeAndUpload(thumbnail)
	if err != nil {
		return nil, err
	}
	return UploaderResult{resizedFid, thumbnailFid}, nil
}

type Uploader struct {
	adapter       StorageAdapter
	maxSize       uint
	thumbnailSize uint
}

func NewUploader(adapter StorageAdapter, maxSize, thumbnailSize uint) ImageUploader {
	return &Uploader{adapter, maxSize, thumbnailSize}
}
