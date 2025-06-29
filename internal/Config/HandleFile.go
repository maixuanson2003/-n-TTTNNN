package Config

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

var url = "http://res.cloudinary.com/dx6b8y6la/video/upload/v1739588778/file-music/file-test.mp4"

func HandleUpLoadFile(input interface{}, publicId string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	env := GetEnvConfig()

	cld, err := cloudinary.NewFromParams(env.GetCloudName(), env.GetCloudApiKey(), env.GetCloudApiSecret())
	if err != nil {
		fmt.Print(err)
		return "", err
	}

	var reader io.Reader

	switch v := input.(type) {
	case *multipart.FileHeader:
		file, err := v.Open()
		if err != nil {
			return "", err
		}
		defer file.Close()
		reader = file
	case multipart.File:
		reader = v
	case io.Reader:
		reader = v
	default:
		return "", fmt.Errorf("unsupported input type: %T", input)
	}

	uploadParam, err := cld.Upload.Upload(ctx, reader, uploader.UploadParams{
		ResourceType: "raw",
		Folder:       env.GetCloudFolder(),
		PublicID:     publicId,
		Format:       "mp3",
	})
	if err != nil {
		fmt.Print(err)
		return "", err
	}

	fmt.Println("Upload thành công:", uploadParam.SecureURL)
	return uploadParam.SecureURL, nil
}
func HandleDownLoadFile(publicId string, types string) (*http.Response, error) {
	url := ""
	switch types {
	case "image":
		url = "http://res.cloudinary.com/dx6b8y6la/%s/upload/v1739588778/file-music/%s.jpg"
	case "video":
		url = "http://res.cloudinary.com/dx6b8y6la/%s/upload/v1739588778/file-music/%s.mp3"
	default:
		return nil, fmt.Errorf("loại file không hợp lệ: %s", types)
	}
	fileUrl := fmt.Sprintf(url, types, publicId)
	resp, err := http.Get(fileUrl)
	fmt.Print(resp)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return resp, nil
}
func HandleUploadImage(input interface{}, publicId string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	env := GetEnvConfig()

	cld, err := cloudinary.NewFromParams(env.GetCloudName(), env.GetCloudApiKey(), env.GetCloudApiSecret())
	if err != nil {
		fmt.Print(err)
		return "", err
	}

	var reader io.Reader

	switch v := input.(type) {
	case *multipart.FileHeader:
		file, err := v.Open()
		if err != nil {
			return "", err
		}
		defer file.Close()
		reader = file
	case multipart.File:
		reader = v
	case io.Reader:
		reader = v
	default:
		return "", fmt.Errorf("unsupported input type: %T", input)
	}

	uploadParam, err := cld.Upload.Upload(ctx, reader, uploader.UploadParams{
		ResourceType: "image",
		Folder:       env.GetCloudFolder(),
		PublicID:     publicId,
	})
	if err != nil {
		fmt.Print(err)
		return "", err
	}

	fmt.Println("Upload ảnh thành công:", uploadParam.SecureURL)
	return uploadParam.SecureURL, nil
}
