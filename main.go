package main

import (
	"flag"
	"fmt"

	"github.com/nishanthvijayan/nptel-dl/nptel"
)

func main() {

	var courseID string
	flag.StringVar(&courseID, "courseID", "", "Course Id or Subject ID (Eg: 106106198)")

	var outputDirectory string
	flag.StringVar(&outputDirectory, "dir", ".", "Output directory")

	var start int
	flag.IntVar(&start, "start", 1, "Video to start at")

	var format string
	flag.StringVar(&format, "format", "mp4", "Options: mp4, 3gp, flv")

	flag.Parse()

	if courseID == "" {
		fmt.Println("No courseID was provided. Exiting..")
		return
	}

	courseIndexPage := nptel.GetCourseIndexPage(courseID)
	defer courseIndexPage.Close()

	courseVideoURLs := nptel.ExtractLectureDownloadUrls(courseIndexPage, format)
	fmt.Printf("Found %d lectures available for download\n", len(courseVideoURLs))

	nptel.DownloadVideos(courseVideoURLs, start, outputDirectory)
}
