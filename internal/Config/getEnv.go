package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
}

func GetEnvConfig() Env {
	return Env{}
}
func (*Env) GetCloudName() string {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
		return ""
	}
	return os.Getenv("CLOUDINARY_CLOUD_NAME")
}
func (*Env) GetCloudApiKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
		return ""
	}
	return os.Getenv("CLOUDINARY_API_KEY")
}
func (*Env) GetCloudFolder() string {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
		return ""
	}
	return os.Getenv("CLOUDINARY_UPLOAD_FOLDER")
}
func (*Env) GetCloudApiSecret() string {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
		return ""
	}
	return os.Getenv("CLOUDINARY_API_SECRET")
}
func (*Env) GetCloudUpLoadPreset() string {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
		return ""
	}
	return os.Getenv("CLOUDINARY_UPLOAD_PRESET")
}
func (*Env) JwtSecretKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
		return ""
	}
	return os.Getenv("JWT_SECRET_KEY")
}
