package utils

import (
	"encoding/base64"
	"os"
	"strings"
)
func UploadBase64Image (base64Image, folder string) string {
	if base64Image == ""{
		return ""
	}
	dataIndex := strings.Index(base64Image, ",")
		if dataIndex == -1 {
			return ""
	}

	data := base64Image[dataIndex+1:]
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return ""
	}

	filename := folder + "/" + generateFilename()
	os.WriteFile(filename, decoded,0644)
	return filename

}

func generateFilename()string {
	return  "image.jpg"
}