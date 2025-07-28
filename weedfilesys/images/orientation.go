package images

import (
	"bytes"
	"github.com/seaweedfs/goexif/exif"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
)

// FixJpgOrientation 根据EXIF中的方向标记修正JPEG图像方向
// data: 输入的JPEG图像字节数据
// 返回: 修正方向后的JPEG图像字节数据
func FixJpgOrientation(data []byte) (oriented []byte) {
	// 解码JPEG图像的EXIF数据
	ex, err := exif.Decode(bytes.NewReader(data))
	if err != nil {
		return data // 如果无法解码EXIF，返回原始数据
	}

	// 获取方向(Orientation)标签
	tag, err := ex.Get(exif.Orientation)
	if err != nil {
		return data // 如果找不到方向标签，返回原始数据
	}

	angle := 0                   // 旋转角度初始化为0
	flipMode := FlipDirection(0) // 翻转模式初始化为不翻转

	// 获取方向标签的整数值
	orient, err := tag.Int(0)
	if err != nil {
		return data // 如果无法获取方向值，返回原始数据
	}

	// 根据不同的方向值设置旋转角度和翻转模式
	switch orient {
	case topLeftSide: // 1: 正常方向，不需要处理
		return data
	case topRightSide: // 2: 水平翻转
		flipMode = 2
	case bottomRightSide: // 3: 旋转180度
		angle = 180
	case bottomLeftSide: // 4: 旋转180度并水平翻转
		angle = 180
		flipMode = 2
	case leftSideTop: // 5: 旋转-90度并水平翻转
		angle = -90
		flipMode = 2
	case rightSideTop: // 6: 旋转-90度
		angle = -90
	case rightSideBottom: // 7: 旋转90度并水平翻转
		angle = 90
		flipMode = 2
	case leftSideBottom: // 8: 旋转90度
		angle = 90
	}

	// 解码原始图像
	if srcImage, _, err := image.Decode(bytes.NewReader(data)); err == nil {
		// 先旋转后翻转
		dstImage := flip(rotate(srcImage, angle), flipMode)

		// 将处理后的图像编码为JPEG格式
		var buf bytes.Buffer
		jpeg.Encode(&buf, dstImage, nil)
		return buf.Bytes()
	}

	return data // 如果图像解码失败，返回原始数据
}

// Exif Orientation Tag 值定义
// 参考: http://sylvana.net/jpegcrop/exif_orientation.html
const (
	topLeftSide     = 1 // 正常方向 (0°)
	topRightSide    = 2 // 水平翻转 (0°镜像)
	bottomRightSide = 3 // 旋转180度 (180°)
	bottomLeftSide  = 4 // 旋转180度并水平翻转 (180°镜像)
	leftSideTop     = 5 // 旋转-90度并水平翻转 (90°CW镜像)
	rightSideTop    = 6 // 旋转-90度 (90°CW)
	rightSideBottom = 7 // 旋转90度并水平翻转 (90°CCW镜像)
	leftSideBottom  = 8 // 旋转90度 (90°CCW)
)

// FlipDirection 类型用于表示图像翻转方向
type FlipDirection int

// 翻转方向常量定义
const (
	FlipVertical   FlipDirection = 1 << iota // 垂直翻转
	FlipHorizontal                           // 水平翻转
)

// DecodeOpts 结构体定义了图像解码选项
type DecodeOpts struct {
	Rotate interface{} // 旋转选项，可以是nil或指定的旋转角度
	Flip   interface{} // 翻转选项，可以是nil或FlipDirection
}

// rotate 根据指定角度旋转图像
// im: 输入图像
// angle: 旋转角度(90, -90, 180, -180)
// 返回: 旋转后的图像
func rotate(im image.Image, angle int) image.Image {
	var rotated *image.NRGBA

	// 根据角度进行不同的旋转处理
	switch angle {
	case 90:
		// 顺时针旋转90度
		newH, newW := im.Bounds().Dx(), im.Bounds().Dy()
		rotated = image.NewNRGBA(image.Rect(0, 0, newW, newH))
		for y := 0; y < newH; y++ {
			for x := 0; x < newW; x++ {
				rotated.Set(x, y, im.At(newH-1-y, x))
			}
		}
	case -90:
		// 逆时针旋转90度
		newH, newW := im.Bounds().Dx(), im.Bounds().Dy()
		rotated = image.NewNRGBA(image.Rect(0, 0, newW, newH))
		for y := 0; y < newH; y++ {
			for x := 0; x < newW; x++ {
				rotated.Set(x, y, im.At(y, newW-1-x))
			}
		}
	case 180, -180:
		// 旋转180度
		newW, newH := im.Bounds().Dx(), im.Bounds().Dy()
		rotated = image.NewNRGBA(image.Rect(0, 0, newW, newH))
		for y := 0; y < newH; y++ {
			for x := 0; x < newW; x++ {
				rotated.Set(x, y, im.At(newW-1-x, newH-1-y))
			}
		}
	default:
		// 不旋转
		return im
	}
	return rotated
}

// flip 根据指定方向翻转图像
// im: 输入图像
// dir: 翻转方向(FlipVertical, FlipHorizontal或其组合)
// 返回: 翻转后的图像
func flip(im image.Image, dir FlipDirection) image.Image {
	if dir == 0 {
		return im // 不需要翻转
	}

	ycbcr := false
	var nrgba image.Image
	dx, dy := im.Bounds().Dx(), im.Bounds().Dy()

	// 检查输入图像是否实现了draw.Image接口
	di, ok := im.(draw.Image)
	if !ok {
		// 如果是YCbCr格式，需要转换为NRGBA
		if _, ok := im.(*image.YCbCr); !ok {
			log.Printf("failed to flip image: input does not satisfy draw.Image")
			return im
		}
		ycbcr = true
		nrgba = image.NewNRGBA(image.Rect(0, 0, dx, dy))
		di, ok = nrgba.(draw.Image)
		if !ok {
			log.Print("failed to flip image: could not cast an NRGBA to a draw.Image")
			return im
		}
	}

	// 水平翻转
	if dir&FlipHorizontal != 0 {
		for y := 0; y < dy; y++ {
			for x := 0; x < dx/2; x++ {
				old := im.At(x, y)
				di.Set(x, y, im.At(dx-1-x, y))
				di.Set(dx-1-x, y, old)
			}
		}
	}

	// 垂直翻转
	if dir&FlipVertical != 0 {
		for y := 0; y < dy/2; y++ {
			for x := 0; x < dx; x++ {
				old := im.At(x, y)
				di.Set(x, y, im.At(x, dy-1-y))
				di.Set(x, dy-1-y, old)
			}
		}
	}

	if ycbcr {
		return nrgba // 返回转换后的NRGBA图像
	}
	return im // 返回翻转后的原始图像
}
