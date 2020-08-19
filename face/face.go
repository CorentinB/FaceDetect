package face

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path"
	"strconv"

	pigo "github.com/esimov/pigo/core"
)

// cropImage takes an image and crops it to the specified rectangle.
func cropImage(img image.Image, crop image.Rectangle) (image.Image, error) {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	// img is an Image interface. This checks if the underlying value has a
	// method called SubImage. If it does, then we can use SubImage to crop the
	// image.
	simg, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}

	return simg.SubImage(crop), nil
}

// Detect detects faces draw rectangles around them.
func Detect(img draw.Image, fileName string, outputDir string) int {
	var facesCount int

	cf, err := Asset("static/facefinder.bin")
	if err != nil {
		log.Fatalf("Error reading the cascade file: %v", err)
	}

	pg := pigo.NewPigo()
	// Unpack the binary file. This will return the number of cascade trees,
	// the tree depth, the threshold and the prediction from tree's leaf nodes.
	classifier, err := pg.Unpack(cf)
	if err != nil {
		log.Fatalf("Error reading the cascade file: %s", err)
	}

	src := pigo.ImgToNRGBA(img)
	pixels := pigo.RgbToGrayscale(src)
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y
	cParams := pigo.CascadeParams{
		MinSize:     20,
		MaxSize:     1000,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,
		ImageParams: pigo.ImageParams{
			Pixels: pixels,
			Rows:   rows,
			Cols:   cols,
			Dim:    cols,
		},
	}

	// Run the classifier over the obtained leaf nodes and return the detection results.
	// The result contains quadruplets representing the row, column, scale and detection score.
	angle := 0.0
	dets := classifier.RunCascade(cParams, angle)

	// Calculate the intersection over union (IoU) of two clusters.
	iouThreshold := 0.2
	faces := classifier.ClusterDetections(dets, iouThreshold)

	// Generate cropped images
	var qThresh float32 = 50.0
	for i, face := range faces {
		if face.Q <= qThresh {
			continue
		}

		cropped, err := cropImage(img, image.Rect(face.Col-face.Scale/2, face.Row-face.Scale/2, face.Col-face.Scale/2+face.Scale, face.Row-face.Scale/2+face.Scale))
		if err != nil {
			log.Fatalf("Error cropping image: %s", err)
		}

		writeImage(cropped, path.Join(outputDir, fileName+"_"+strconv.Itoa(i)+".png"))

		facesCount++
	}
	return facesCount
}

// writeImage writes an Image back to the disk.
func writeImage(img image.Image, name string) error {
	fd, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fd.Close()

	return png.Encode(fd, img)
}
