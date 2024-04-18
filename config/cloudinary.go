package config

import (
	"github.com/cloudinary/cloudinary-go/v2"
)

func SetupCloudinary() (*cloudinary.Cloudinary, error) {
	cldSecret := "uFJhxNAJlCeLqQbSqp8-pV3H7lE"
	cldName := "dytequc2m"
	cldKey := "283248446939275"

	cld, err := cloudinary.NewFromParams(cldName, cldKey, cldSecret)
	if err != nil {
		return nil, err
	}

	return cld, nil
}
