package Config

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

var url = "http://res.cloudinary.com/dx6b8y6la/video/upload/v1739588778/file-music/file-test.mp4"

func HandleUpLoadFile(input interface{}, publicId string) (string, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 60*time.Minute)
	Env := GetEnvConfig()

	defer cancle()
	cld, errorToUpLoadFile := cloudinary.NewFromParams(Env.GetCloudName(), Env.GetCloudApiKey(), Env.GetCloudApiSecret())
	if errorToUpLoadFile != nil {
		fmt.Print(errorToUpLoadFile)
		return "", errorToUpLoadFile
	}
	UploadParam, err := cld.Upload.Upload(ctx, input, uploader.UploadParams{
		Folder:       Env.GetCloudFolder(),
		UploadPreset: Env.GetCloudUpLoadPreset(),
		PublicID:     publicId,
	})
	fmt.Print(UploadParam)
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	fmt.Println("Upload thành công:", UploadParam.SecureURL)

	return UploadParam.SecureURL, nil
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
