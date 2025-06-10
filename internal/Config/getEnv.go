package Config

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

func (*Env) SmtpHost() string {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
		return ""
	}
	return os.Getenv("SMTP_HOST")
}

func (*Env) SmtpPort() string {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
		return ""
	}
	return os.Getenv("SMTP_PORT")
}

func (*Env) SmtpUser() string {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
		return ""
	}
	return os.Getenv("SMTP_USER")
}

func (*Env) SmtpPassword() string {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
		return ""
	}
	return os.Getenv("SMTP_PASSWORD")
}

func (*Env) FromEmail() string {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
		return ""
	}
	return os.Getenv("FROM_EMAIL")
}

func (*Env) GeminiAiKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
		return ""
	}
	return os.Getenv("GEMINI_API_KEY")
}
