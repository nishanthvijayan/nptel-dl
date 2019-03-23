package nptel

import (
	"fmt"
	"strings"
)

type lectureVideo struct {
	topicName   string
	videoID     string
	format      string
	downloadURL string
}

func newLectureVideo(downloadURL string) *lectureVideo {
	urlParts := strings.Split(downloadURL, "=")
	videoID := strings.TrimSuffix(strings.Split(urlParts[2], "&")[0], ".mp4")
	topicName := urlParts[len(urlParts)-1]

	return &lectureVideo{
		topicName:   topicName,
		videoID:     videoID,
		downloadURL: downloadURL,
	}
}

func (lecture lectureVideo) generateFileName(outputDirectory string) string {
	return fmt.Sprintf("%s/%s-%s.mp4",
		strings.TrimRight(outputDirectory, "/"),
		lecture.videoID,
		lecture.topicName,
	)
}
