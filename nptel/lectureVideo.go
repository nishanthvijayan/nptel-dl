package nptel

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

type lectureVideo struct {
	topicName   string
	videoID     string
	format      string
	downloadURL string
}

func newLectureVideo(downloadURL string) *lectureVideo {
	queryParams, err := url.ParseQuery(downloadURL)
	if err != nil {
		log.Fatal("Malformed url found. Aborting")
	}

	filename := queryParams.Get("filename")
	videoID := strings.Split(filename, ".")[0]
	format := strings.Split(filename, ".")[1]
	topicName := queryParams.Get("subjectName")

	return &lectureVideo{
		topicName:   topicName,
		videoID:     videoID,
		downloadURL: downloadURL,
		format:      format,
	}
}

func (lecture lectureVideo) generateFileName(outputDirectory string) string {
	return fmt.Sprintf("%s/%s-%s.%s",
		strings.TrimRight(outputDirectory, "/"),
		lecture.videoID,
		lecture.topicName,
		lecture.format,
	)
}
