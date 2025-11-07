package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
)

// ReadAsBytes return data, mimeType, error
func ReadAsBytes(fullFilePath string) ([]byte, string, error) {
	bytes, err := ioutil.ReadFile(fullFilePath)
	if err != nil {
		return nil, "", err
	}

	mimeType := http.DetectContentType(bytes)

	return bytes, mimeType, nil
}

// ReadImage return Image, error
func ReadImage(fileFullPath string) (image.Image, error) {
	return imaging.Open(fileFullPath)
}

// ReadAsBase64 return data, error
func ReadAsBase64(fullFilePath string) (string, error) {
	bytes, err := ioutil.ReadFile(fullFilePath)
	if err != nil {
		return "", err
	}

	var base64Encoding string

	mimeType := http.DetectContentType(bytes)

	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	case "image/bmp":
		base64Encoding += "data:image/bmp;base64,"
	}

	base64Encoding += ByteArrayToBase64(bytes)

	return base64Encoding, nil
}

// WriteImage func, write image to fileFullPath, return an error
func WriteImage(img image.Image, fileFullPath string) error {
	return imaging.Save(img, fileFullPath)
}

// ByteArrayToBase64 func, convert array of bytes to base64 string
func ByteArrayToBase64(b []byte) string {
	img, err := imaging.Decode(bytes.NewReader(b), imaging.AutoOrientation(true))
	if err != nil {
		return base64.StdEncoding.EncodeToString(b)
	}
	var buf bytes.Buffer
	// detect content type
	contentType := http.DetectContentType(b)
	switch contentType {
	case "image/png":
		err = png.Encode(&buf, img)
	case "image/jpeg":
		err = jpeg.Encode(&buf, img, nil)
	default:
		return base64.StdEncoding.EncodeToString(b)
	}
	if err != nil {
		return base64.StdEncoding.EncodeToString(b)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

// Base64ToByteArray func, convert base64 string to array of bytes
func Base64ToByteArray(str string) ([]byte, error) {
	parts := strings.Split(str, ",")
	return base64.StdEncoding.DecodeString(parts[len(parts)-1])
}

// ByteArrayToImage func, convert array of bytes to an image
func ByteArrayToImage(data []byte) (image.Image, string, error) {
	return image.Decode(bytes.NewReader(data))
}

// ImageToByteArray func, convert an image to bytes array
func ImageToByteArray(img image.Image, imageType string) ([]byte, error) {
	buf := new(bytes.Buffer)
	if imageType == "image/jpeg" {
		if err := jpeg.Encode(buf, img, nil); err != nil {
			return nil, err
		}
	} else if imageType == "image/png" {
		if err := png.Encode(buf, img); err != nil {
			return nil, err
		}
	} else if imageType == "image/bmp" {
		if err := png.Encode(buf, img); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Unsupported image type")
	}
	return buf.Bytes(), nil
}

// ResizeImageWithBytes func, resize bytes array image
// width = 0 or height = 0 meaning keep the ratio
func ResizeImageWithBytes(data []byte, width, height int) ([]byte, error) {
	sourceImg, imageType, err := ByteArrayToImage(data)
	if err != nil {
		return nil, err
	}
	destImg := imaging.Resize(sourceImg, width, height, imaging.Lanczos)
	return ImageToByteArray(destImg, "image/"+imageType)
}

// ResizeImageWithBase64 func, resize base64 image
// width = 0 or height = 0 meaning keep the ratio
func ResizeImageWithBase64(data string, width, height int) (string, error) {
	dataBytes, err := Base64ToByteArray(data)
	if err != nil {
		return "", err
	}

	destImgAsBytes, err := ResizeImageWithBytes(dataBytes, width, height)

	if err != nil {
		return "", err
	}

	return ByteArrayToBase64(destImgAsBytes), nil
}

// ResizeImageWithBase64 func, resize an image
// width = 0 or height = 0 meaning keep the ratio
func ResizeImage(img image.Image, width, height int) image.Image {
	return imaging.Resize(img, width, height, imaging.Lanczos)
}

// CropImage func, crop an image from the center
func CropImage(img image.Image, width, height int) (image.Image, error) {
	return imaging.CropCenter(img, width, height), nil
}

func cropImageWithBytes(data []byte, x0, y0, x1, y1 int) ([]byte, error) {
	sourceImg, imageType, err := ByteArrayToImage(data)
	if err != nil {
		return nil, err
	}

	destImg := imaging.Crop(sourceImg, image.Rect(x0, y0, x1, y1))
	return ImageToByteArray(destImg, "image/"+imageType)
}

// CropImageWithBytes func, crop bytes array image from the center
func CropImageWithBytes(data []byte, width, height int) ([]byte, error) {
	sourceImg, imageType, err := ByteArrayToImage(data)
	if err != nil {
		return nil, err
	}

	destImg := imaging.CropCenter(sourceImg, width, height)
	return ImageToByteArray(destImg, "image/"+imageType)
}

// CropImageWithBase64 func, crop base64 image from the center
func CropImageWithBase64(data string, width, height int) (string, error) {
	dataBytes, err := Base64ToByteArray(data)
	if err != nil {
		return "", err
	}

	destImgAsBytes, err := CropImageWithBytes(dataBytes, width, height)

	if err != nil {
		return "", err
	}

	return ByteArrayToBase64(destImgAsBytes), nil
}

// CropImageWithBytesByCoordinates func, crop bytes array image from left, top, width, height
func CropImageWithBytesByCoordinates(data []byte, left, top, width, height int) ([]byte, error) {
	return cropImageWithBytes(data, left, top, width, height)
}

// CropImageWithBase64ByCoordinates func, crop base64 image from left, top, width, height
func CropImageWithBase64ByCoordinates(data string, left, top, width, height int) (string, error) {
	dataBytes, err := Base64ToByteArray(data)
	if err != nil {
		return "", err
	}
	destImgAsBytes, err := cropImageWithBytes(dataBytes, left, top, width, height)

	if err != nil {
		return "", err
	}
	return ByteArrayToBase64(destImgAsBytes), nil
}

// Convert top, left, right, bottom witdh image size to (x0, y0) (x1, y1) coordinates
func convertTLRBPercentWithImageSizetoPoint(top, left, right, bottom, width, height int) (x0, y0, x1, y1 int) {
	x0, y0, x1, y1 = 0, 0, 0, 0

	x0 = (left * width) / 100
	y0 = (top * height) / 100

	x1 = ((100 - right) * width) / 100
	y1 = ((100 - bottom) * height) / 100

	return x0, y0, x1, y1
}

// CropImageWithBase64ByMargin func, crop base64 image from left, top, right, bottom
func CropImageByMargin(data string, top, left, right, bottom int) (string, error) {
	if data == "" {
		return "", errors.New("Base64 empty")
	}

	if left == 0 && top == 0 && right == 0 && bottom == 0 {
		return data, nil
	}

	parts := strings.Split(data, ",")
	dataBytes, err := Base64ToByteArray(parts[len(parts)-1])
	if err != nil {
		return "", err
	}

	image, _, err := image.DecodeConfig(bytes.NewReader(dataBytes))
	if err != nil {
		fmt.Println(err)
	}
	width, height := image.Width, image.Height

	x0, y0, x1, y1 := convertTLRBPercentWithImageSizetoPoint(top, left, right, bottom, width, height)

	destImgAsBytes, err := cropImageWithBytes(dataBytes, x0, y0, x1, y1)

	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v,%v", parts[0], ByteArrayToBase64(destImgAsBytes)), nil
}
