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

	plusSign := strings.HasPrefix(ph, "+") || strings.HasPrefix(ph, "00")

	ph = convertArabicToLatin(ph)
	ph = regexp.MustCompile(`\D`).ReplaceAllString(ph, "")
	ph = strings.TrimLeft(ph, "0")

	if plusSign {
		ph = "+" + ph
	}

	return ph, nil
}
