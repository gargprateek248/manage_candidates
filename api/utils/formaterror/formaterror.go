package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {

	if strings.Contains(err, "phone_no") {
		return errors.New("Phone Number Already Taken")
	}

	return errors.New("Incorrect Details")
}
