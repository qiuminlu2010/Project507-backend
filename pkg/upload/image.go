package upload

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"

	"qiu/blog/pkg/file"
	"qiu/blog/pkg/setting"
	"qiu/blog/pkg/util"

	"github.com/disintegration/imaging"
)

func GetImageFullUrl(name string) string {
	return setting.AppSetting.ImagePrefixUrl + "/" + GetImagePath() + name
}

func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName + time.Now().String())

	return fileName + ext
}

func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

func GetImageTempPath() string {
	return setting.AppSetting.ImageTempSavePath
}

func GetThumbPath() string {
	return setting.AppSetting.ThumbSavePath
}
func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.EqualFold(allowExt, ext) {
			return true
		}
	}

	return false
}

func CheckImageSize(f *multipart.FileHeader) bool {
	// size, err := file.GetSize(f)
	// if err != nil {
	// 	log.Println(err)
	// 	logging.Warn(err)
	// 	return false
	// }
	// logging.Info("上传图片大小", size, setting.AppSetting.ImageMaxSize)
	return int(f.Size) <= setting.AppSetting.ImageMaxSize
}

func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := file.CheckPermission(src)
	if perm {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}

func Thumbnailify(fileName string) (outputPath string, err error) {
	// 读取文件
	var (
		file *os.File
		img  image.Image
	)
	imagePath := GetImagePath() + fileName
	if file, err = os.Open(imagePath); err != nil {
		return
	}
	defer file.Close()
	extname := strings.ToLower(path.Ext(imagePath))
	outputPath = GetThumbPath() + fileName
	// decode jpeg into image.Image
	switch extname {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	case ".gif":
		img, err = gif.Decode(file)
	default:
		err = errors.New("Unsupport file type" + extname)
		return
	}

	if img == nil {
		err = errors.New("generate thumbnail fail")
		return
	}
	thumb := imaging.Resize(img, setting.AppSetting.ThumbMaxWidth, 0, imaging.Lanczos)
	// thumb := thumbnailCrop(512, 512, img)
	// thumb := resize.Thumbnail(300, 600, img, resize.NearestNeighbor)
	out, err := os.Create(outputPath)
	if err != nil {
		return
	}
	defer out.Close()
	switch extname {
	case ".jpg", ".jpeg":
		jpeg.Encode(out, thumb, nil)
	case ".png":
		png.Encode(out, thumb)
	case ".gif":
		gif.Encode(out, thumb, nil)
	default:
		err = errors.New("Unsupport file type" + extname)
		return
	}
	return
}

// 缩略图按照指定的宽和高非失真缩放裁剪
// func thumbnailCrop(minWidth, minHeight uint, img image.Image) image.Image {
// 	origBounds := img.Bounds()
// 	origWidth := uint(origBounds.Dx())
// 	origHeight := uint(origBounds.Dy())
// 	newWidth, newHeight := origWidth, origHeight

// 	// Return original image if it have same or smaller size as constraints
// 	if minWidth >= origWidth && minHeight >= origHeight {
// 		return img
// 	}

// 	if minWidth > origWidth {
// 		minWidth = origWidth
// 	}

// 	if minHeight > origHeight {
// 		minHeight = origHeight
// 	}

// 	// Preserve aspect ratio
// 	if origWidth > minWidth {
// 		newHeight = uint(origHeight * minWidth / origWidth)
// 		if newHeight < 1 {
// 			newHeight = 1
// 		}
// 		newWidth = minWidth
// 	}

// 	if newHeight < minHeight {
// 		newWidth = uint(newWidth * minHeight / newHeight)
// 		if newWidth < 1 {
// 			newWidth = 1
// 		}
// 		newHeight = minHeight
// 	}

// 	thumbImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
// 	return thumbImg
// }
