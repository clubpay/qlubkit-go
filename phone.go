package qkit

import (
	"regexp"
	"strings"
)

func convertArabicToLatin(input string) string {
	replacer := strings.NewReplacer("٠", "0", "١", "1", "٢", "2", "٣", "3", "٤", "4", "٥", "5", "٦", "6", "٧", "7", "٨", "8", "٩", "9")
	out := replacer.Replace(input)

	return out
}

// SanitizePhoneNumber sanitizes given phone number
func SanitizePhoneNumber(phone string) (string, error) {
	ph := phone

	ph = convertArabicToLatin(ph)

	re1, err := regexp.Compile("[^0-9]")
	if err != nil {
		return "", err
	}

	ph = re1.ReplaceAllString(ph, "")

	re2, err := regexp.Compile("^0+")
	if err != nil {
		return "", err
	}

	ph = re2.ReplaceAllString(ph, "")

	ph = "+" + ph

	return ph, nil
}
