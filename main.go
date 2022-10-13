package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/echo/v4"

	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
)

const (
	AWS_S3_REGION = "ap-northeast-2"
	AWS_S3_BUCKET = "ssh-dump-files"
)

func main() {

	fmt.Println("hi")

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "health",
		})
	})

	e.GET("/:key", func(c echo.Context) error {

		// 워터마크 파일 추가
		wmb, _ := os.Open("gsn.png")
		watermark, _ := png.Decode(wmb)
		defer wmb.Close()

		bucket := AWS_S3_BUCKET
		// S3 경로 변경
		key := "iu/" + c.Param("key")

		session, err := session.NewSession(&aws.Config{Region: aws.String(AWS_S3_REGION)})
		if err != nil {
			fmt.Println(err)
		}

		svc := s3.New(session)

		result, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})
		if err != nil {
			fmt.Printf("Failed to get data to %s/%s, %s\n", bucket, key, err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": c.Param("key"),
			})
		}
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(result.Body)

		fmt.Println("Image Process")

		img, _ := jpeg.Decode(buf)

		offset := image.Pt(20, 20)

		b := img.Bounds()
		m := image.NewRGBA(b)
		draw.Draw(m, b, img, image.ZP, draw.Src)
		draw.Draw(m, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)

		newBuf := new(bytes.Buffer)
		jpeg.Encode(newBuf, m, &jpeg.Options{jpeg.DefaultQuality})

		// return c.Blob(200, "image/jpg", buf.Bytes())
		return c.Blob(200, "image/jpg", newBuf.Bytes())
	})

	e.Logger.Debug(e.Start(":8080"))
}
