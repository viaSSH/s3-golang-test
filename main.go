package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/echo/v4"
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

		bucket := "ssh-dump-files"
		key := "iu/" + c.Param("key")

		session, err := session.NewSession(&aws.Config{Region: aws.String(AWS_S3_REGION)})
		if err != nil {
			// log.Fatal(err)
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

		return c.Blob(200, "image/jpg", buf.Bytes())
	})

	e.Logger.Debug(e.Start(":8080"))
}
