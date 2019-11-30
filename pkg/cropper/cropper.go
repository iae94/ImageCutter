package cropper

import (
	cfg "ImageCutter/pkg/config"
	"ImageCutter/pkg/models"
	"bytes"
	"github.com/disintegration/imaging"
	"go.uber.org/zap"
	"golang.org/x/image/tiff"
	"image/gif"
	"image/jpeg"
	"image/png"
	"path/filepath"
)


type Cropper struct {
	Logger *zap.Logger
	Config *cfg.CutterConfig
}

func NewCropper(logger *zap.Logger, config *cfg.CutterConfig) *Cropper {
	return &Cropper{
		Logger: logger,
		Config: config,
	}
}



func (c *Cropper) Crop(width int, height int, image *models.Image) ([]byte, error){

	imagePath := filepath.Join(c.Config.Cutter.Cache.Folder, image.Name)

	img, err := imaging.Open(imagePath)
	if err != nil {
		c.Logger.Sugar().Errorf("Cropper cannot open image: %v error: %v", imagePath, err)
		return nil, err
	}

	resizedImage := imaging.Resize(img, width, height, imaging.Lanczos)

	buffer := new(bytes.Buffer)
	switch image.MimeType {
	case "image/jpeg":
		err = jpeg.Encode(buffer, resizedImage, nil)
	case "image/png":
		err = png.Encode(buffer, resizedImage)
	case "image/tiff":
		err = tiff.Encode(buffer, resizedImage, nil)
	case "image/gif":
		err = gif.Encode(buffer, resizedImage, nil)
	default:
		err = jpeg.Encode(buffer, resizedImage, nil)
	}

	if err != nil {
		c.Logger.Sugar().Errorf("Cropper cannot convert image.NRGBA of %v to []byte | Error: %v", image.MimeType, err)
		return nil, err
	}
	croppedImage := buffer.Bytes()

	return croppedImage, nil
}
