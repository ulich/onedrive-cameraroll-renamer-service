package internal

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func calcNewFilename(name string) (string, error) {
	if matches(`^20\d{6}_\d{6}\.(jpg|jpeg|mp4)$`, name) ||
		matches(`^20\d{6}_\d{6}_\d+\.(jpg|jpeg|mp4)$`, name) ||
		matches(`^20\d{6}_\d{6}\(\d+\)\.(jpg|jpeg|mp4)$`, name) ||
		matches(`^20\d{6}_000000_\d+_WA\.(jpg|jpeg|mp4)$`, name) {
		return name, nil
	}

	if matches(`^20\d{6}_\d{6}\.(JPG|JPEG|MP4)$`, name) ||
		matches(`^20\d{6}_\d{6}_\d+\.(JPG|JPEG|MP4)$`, name) {
		return strings.ToLower(name), nil
	}

	if matches(`^IMG_20\d{6}_\d{6}\.(jpg|jpeg)$`, name) ||
		matches(`^IMG_20\d{6}_\d{6}_\d+\.(jpg|jpeg)$`, name) {
		return strings.ReplaceAll(name, "IMG_", ""), nil
	}

	if matches(`^IMG_\d{13}_\d+\.(jpg|jpeg)$`, name) ||
		matches(`^IMAGE_\d{13}_\d+\.(jpg|jpeg)$`, name) {

		parts := strings.Split(name, "_")
		timestampString := parts[1]
		suffix := parts[2]

		timestamp, err := strconv.ParseInt(timestampString, 10, 64)
		if err != nil {
			return "", fmt.Errorf("parse timestamp as int %v, %w", timestampString, err)
		}

		timezone, err := time.LoadLocation("Europe/Berlin")
		if err != nil {
			return "", fmt.Errorf("load timezone: %w", err)
		}

		t := time.UnixMilli(timestamp)
		timeString := t.In(timezone).Format("20060102_150405")

		return timeString + "_" + suffix, nil
	}

	if matches(`^IMG-20\d{6}-WA\d{4}\.(jpg|jpeg)$`, name) {
		newName := strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(name, "IMG-", ""),
					"WA", ""),
				"-", "_"),
			".", "_WA.")
		return newName, nil
	}

	if matches(`^VID_20\d{6}_\d{6}\.mp4$`, name) {
		return strings.ReplaceAll(name, "VID_", ""), nil
	}

	if matches(`^VID-20\d{6}-WA\d{4}\.mp4$`, name) {
		newName := strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(name, "VID-", ""),
					"WA", ""),
				"-", "_"),
			".", "_WA.")
		return newName, nil
	}

	return "", fmt.Errorf("unknown filename pattern %s", name)
}

func matches(pattern, name string) bool {
	exp := regexp.MustCompile(pattern)
	return exp.MatchString(name)
}
