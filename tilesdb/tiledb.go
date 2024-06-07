package tilesdb

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/kyroy/kdtree"
	"github.com/kyroy/kdtree/points"
)

type TileDB struct {
	kdtree.KDTree
}

var DefaultTilesDB *TileDB = NewTileDB()

func NewTileDB() *TileDB {
	var db TileDB
	db.Init()
	return &db
}

// 获取整张图片的平均rgb值
func AverageColor(img image.Image) []float64 {
	bounds := img.Bounds()
	rsum, gsum, bsum := 0.0, 0.0, 0.0
	//遍历图片所有点，把每个点的rgb值累加起来
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rsum, gsum, bsum = rsum+float64(r), gsum+float64(g), bsum+float64(b)
		}
	}
	totalPixels := float64(bounds.Dx() * bounds.Dy())
	ret := make([]float64, 0, 3)
	ret = append(ret, rsum/totalPixels, gsum/totalPixels, bsum/totalPixels)
	return ret
}

// 把图片缩放到指定的尺寸
func Resize(img image.Image, newWidth int) image.NRGBA {
	bounds := img.Bounds()
	ratio := bounds.Dx() / newWidth
	out := image.NewNRGBA(image.Rect(bounds.Min.X/ratio, bounds.Min.Y/ratio,
		bounds.Max.X/ratio, bounds.Max.Y/ratio))
	for y, j := bounds.Min.Y, bounds.Min.Y; y < bounds.Max.Y; y, j = y+ratio, j+1 {
		for x, i := bounds.Min.X, bounds.Min.X; x < bounds.Max.X; x, i = x+ratio, i+1 {
			r, g, b, a := img.At(x, y).RGBA()
			out.SetNRGBA(i, j, color.NRGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
		}
	}

	return *out
}

// 初始化，从tiles中加载图片
func (db *TileDB) Init() {
	var dirname string = "tiles"
	fmt.Println("开始构建嵌入图片数据库...")
	files, _ := os.ReadDir(dirname)
	for _, f := range files {
		name := "tiles/" + f.Name()
		file, err := os.Open(name)
		if err == nil {
			img, _, err := image.Decode(file)
			if err == nil {
				avgColor := AverageColor(img)
				pic := points.NewPoint(avgColor, name)
				db.Insert(pic)
			} else {
				fmt.Println("构建嵌入图片数据库出错：", err, name)
			}
		} else {
			fmt.Println("构建嵌入图片数据库出错：", err, "无法打开文件", name)
		}
		file.Close()
	}
	fmt.Println("构建嵌入图片数据库完毕")
}

func (db *TileDB) Nearest(color []float64) string {
	target := db.KNN(points.NewPoint(color, ""), 1)
	if pic, ok := target[0].(*points.Point); ok {
		return pic.Data.(string)
	} else {
		return "tiles/0A1E54CDCA89EA5C420162ADE3E96FAF.jpg" // 找不到图片时返回默认图片
	}
}

func CloneTilesDB() *TileDB {
	return DefaultTilesDB
}
