package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "https://nptel.ac.in"

type lectureVideo struct {
	topicName   string
	videoID     string
	downloadURL string
}

func NewLectureVideoFromDownloadURL(downloadURL string) lectureVideo {
	urlParts := strings.Split(downloadURL, "=")
	videoID := strings.TrimSuffix(strings.Split(urlParts[2], "&")[0], ".mp4")
	topicName := urlParts[len(urlParts)-1]

	return lectureVideo{
		topicName:   topicName,
		videoID:     videoID,
		downloadURL: downloadURL,
	}
}

func createFileNameFromLectureVideo(lecture lectureVideo, outputDirectory string) string {
	return fmt.Sprintf("%s/%s-%s.mp4",
		strings.TrimRight(outputDirectory, "/"),
		lecture.videoID,
		lecture.topicName,
	)
}

func escapeSpaceInURL(url string) string {
	return strings.Replace(url, " ", "%20", -1)
}

func getCourseIndexPage(courseID string) io.ReadCloser {

	res, err := http.Get("https://nptel.ac.in/courses/nptel_download.php?subjectid=" + courseID)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return res.Body
}

func isMp4DownloadLink(s *goquery.Selection) bool {
	return strings.Contains(s.Text(), "MP4 Download")
}

func extractLectureDownloadUrls(page io.Reader) []string {

	doc, err := goquery.NewDocumentFromReader(page)
	if err != nil {
		log.Fatal(err)
	}

	urls := []string{}
	const downloadLinkSelector = "td > a"
	doc.Find(downloadLinkSelector).Each(func(i int, s *goquery.Selection) {
		if isMp4DownloadLink(s) {
			if relativeDownloadURL, exists := s.Attr("href"); exists {
				urls = append(urls, baseURL+relativeDownloadURL)
			}
		}
	})

	return urls
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(escapeSpaceInURL(url))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func downloadVideos(videoURLs []string, start int, outputDirectory string) {
	for i, videoURL := range videoURLs {

		lectureVideo := NewLectureVideoFromDownloadURL(videoURL)

		// i+1 because i is zero-indexed whereas start is 1-indexed
		if i+1 < start {
			fmt.Printf("Skipping Video %d  - %s\n", i+1, lectureVideo.topicName)
			continue
		}

		fileName := createFileNameFromLectureVideo(lectureVideo, outputDirectory)
		fmt.Printf("Downloading Video %d  - %s\n", i+1, lectureVideo.topicName)
		fmt.Printf("Saving it to: %s\n\n", fileName)

		err := downloadFile(fileName, videoURL)
		if err != nil {
			panic(err)
		}
	}
}

func main() {

	var courseID string
	flag.StringVar(&courseID, "courseID", "", "Course Id or Subject ID (Eg: 106106198)")

	var outputDirectory string
	flag.StringVar(&outputDirectory, "dir", ".", "Output directory")

	var start int
	flag.IntVar(&start, "start", 1, "Video to start at (default is 1)")

	flag.Parse()

	if courseID == "" {
		fmt.Println("No courseID was provided. Exiting..")
		return
	}

	courseIndexPage := getCourseIndexPage(courseID)
	defer courseIndexPage.Close()

	courseVideoURLs := extractLectureDownloadUrls(courseIndexPage)

	downloadVideos(courseVideoURLs, start, outputDirectory)
}
