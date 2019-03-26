package nptel

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	pb "gopkg.in/cheggaaa/pb.v1"
)

const baseURL = "https://nptel.ac.in"

// GetCourseIndexPage takes as input a courseID and returns the body
// of the course index page of the corresponding course
// If page was not found or unreachable, the function will exist after logging the status code
func GetCourseIndexPage(courseID string) io.ReadCloser {

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

func isFLVDownloadLink(s *goquery.Selection) bool {
	return strings.Contains(s.Text(), "FLV Download")
}

func is3GPDownloadLink(s *goquery.Selection) bool {
	return strings.Contains(s.Text(), "3GP Download")
}

func getDownloadLinkSelector(format string) func(s *goquery.Selection) bool {
	switch format {
	case "mp4":
		return isMp4DownloadLink
	case "flv":
		return isFLVDownloadLink
	case "3gp":
		return is3GPDownloadLink
	default:
		return isMp4DownloadLink
	}
}

// ExtractLectureDownloadUrls reads the contents of a Reader and extract
// urls that look like MP4 download links
func ExtractLectureDownloadUrls(page io.Reader, format string) []string {
	isDownloadLink := getDownloadLinkSelector(format)

	doc, err := goquery.NewDocumentFromReader(page)
	if err != nil {
		log.Fatal(err)
	}

	urls := []string{}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if isDownloadLink(s) {
			if relativeDownloadURL, exists := s.Attr("href"); exists {
				urls = append(urls, baseURL+relativeDownloadURL)
			}
		}
	})

	return urls
}

func ioCopyWithReporting(r io.Reader, w io.Writer, totalBytes int) error {
	bar := pb.New(totalBytes).SetUnits(pb.U_BYTES)

	bar.Format(" \u2580-- ")
	bar.ShowSpeed = true

	bar.Start()

	barProxyReader := bar.NewProxyReader(r)
	_, err := io.Copy(w, barProxyReader)

	bar.Finish()

	return err
}

// downloadFile saves the contents of url to destinationFilePath
func downloadFile(url string, destinationFilepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fileSize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	out, err := os.Create(destinationFilepath)
	if err != nil {
		return err
	}
	defer out.Close()

	return ioCopyWithReporting(resp.Body, out, fileSize)
}

func escapeSpaceInURL(url string) string {
	return strings.Replace(url, " ", "%20", -1)
}

// DownloadVideos downloads the videos using the urls provided in videoUrls
// If start is specified to greater that 1, then videos before the start index will not be downloaded
// saves them to the directory provided in outputDirectory
func DownloadVideos(videoURLs []string, start int, outputDirectory string) {
	for i, videoURL := range videoURLs {

		lectureVideo := newLectureVideo(videoURL)

		// i+1 because i is zero-indexed whereas start is 1-indexed
		if i+1 < start {
			fmt.Printf("Skipping Video %d  - %s\n", i+1, lectureVideo.topicName)
			continue
		}

		fileName := lectureVideo.generateFileName(outputDirectory)
		fmt.Printf("Downloading Video %d  - %s\n", i+1, lectureVideo.topicName)
		fmt.Printf("Saving it to: %s\n\n", fileName)

		err := downloadFile(escapeSpaceInURL(videoURL), fileName)
		if err != nil {
			panic(err)
		}
	}
}
