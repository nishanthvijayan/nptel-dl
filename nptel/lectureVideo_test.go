package nptel

import "testing"

func Test_NewLectureVideoParsesMp4Link(t *testing.T) {
	mp4DownloadURL := "https://nptel.ac.in/courses/download_mp4.php?subjectId=106106198&filename=mod01lec05.mp4&subjectName=Norms"
	lectureVideo := newLectureVideo(mp4DownloadURL)

	if lectureVideo.videoID != "mod01lec05" {
		t.Errorf("VideoID was not parsed properly")
	}

	if lectureVideo.format != "mp4" {
		t.Errorf("Format was not parsed properly")
	}

	if lectureVideo.topicName != "Norms" {
		t.Errorf("Topic Name was not parsed properly")
	}
}

func Test_NewLectureVideoParses3GPLink(t *testing.T) {
	mp4DownloadURL := "https://nptel.ac.in/courses/download_3gp.php?subjectId=106106198&filename=mod01lec04.3gp&subjectName=Basic%20Operations"
	lectureVideo := newLectureVideo(mp4DownloadURL)

	if lectureVideo.videoID != "mod01lec04" {
		t.Errorf("VideoID was not parsed properly")
	}

	if lectureVideo.format != "3GP" {
		t.Errorf("Format was not parsed properly")
	}

	if lectureVideo.topicName != "Basic%20Operations" {
		t.Errorf("VideoID was not parsed properly")
	}
}

func Test_NewLectureVideoParsesFLVLink(t *testing.T) {
	mp4DownloadURL := "https://nptel.ac.in/courses/download_flv.php?subjectId=106106198&filename=mod02lec12.flv&subjectName=Expectation"
	lectureVideo := newLectureVideo(mp4DownloadURL)

	if lectureVideo.videoID != "mod02lec12" {
		t.Errorf("VideoID was not parsed properly")
	}

	if lectureVideo.format != "FLV" {
		t.Errorf("Format was not parsed properly")
	}

	if lectureVideo.topicName != "Expectation" {
		t.Errorf("VideoID was not parsed properly")
	}
}
