package Utils

import (
	"encoding/base64"
	"fmt"
	"math/rand"

)

func Token(length int) (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "err", fmt.Errorf("rastgele UID üretilemedi: %v", err)
	}

	return base64.URLEncoding.EncodeToString(bytes), nil

}