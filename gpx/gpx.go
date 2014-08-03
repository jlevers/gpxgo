package gpx

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"time"
	//"fmt"

	"github.com/tkrajina/gpxgo/gpx/gpx10"
	"github.com/tkrajina/gpxgo/gpx/gpx11"

	//"fmt"
)

// An array cannot be constant :( The first one if the default layout:
var TIMELAYOUTS = []string{
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05Z",
	"2006-01-02 15:04:05",
}

type ToXmlParams struct {
	Version string
	Indent  bool
}

/*
 * Params are optional, you can set null to use GPXs Version and no indentation.
 */
func ToXml(g *GPX, params ToXmlParams) ([]byte, error) {
	version := g.Version
	if len(params.Version) > 0 {
		version = params.Version
	}
	indentation := params.Indent

	var gpxDoc interface{}
	if version == "1.0" {
		gpxDoc = convertToGpx10Models(g)
	} else if version == "1.1" {
		gpxDoc = convertToGpx11Models(g)
	} else {
		return nil, errors.New("Invalid version " + version)
	}

	var buffer bytes.Buffer
	buffer.WriteString(xml.Header)
	if indentation {
		bytes, err := xml.MarshalIndent(gpxDoc, "", "	")
		if err != nil {
			return nil, err
		}
		buffer.Write(bytes)
	} else {
		bytes, err := xml.Marshal(gpxDoc)
		if err != nil {
			return nil, err
		}
		buffer.Write(bytes)
	}
	return buffer.Bytes(), nil
}

func guessGPXVersion(bytes []byte) (string, error) {
	startOfDocument := string(bytes[:1000])

	parts := strings.Split(startOfDocument, "<gpx")
	if len(parts) <= 1 {
		return "", errors.New("Invalid GPX file, cannot find version")
	}
	parts = strings.Split(parts[1], "version=")

	if len(parts) <= 1 {
		return "", errors.New("Invalid GPX file, cannot find version")
	}

	if len(parts[1]) < 10 {
		return "", errors.New("Invalid GPX file, cannot find version")
	}

	result := parts[1][1:4]

	return result, nil
}

func parseGPXTime(timestr string) (*time.Time, error) {
	if strings.Contains(timestr, ".") {
		// Probably seconds with milliseconds
		timestr = strings.Split(timestr, ".")[0]
	}
	timestr = strings.Trim(timestr, " \t\n\r")
	for _, timeLayout := range TIMELAYOUTS {
		t, err := time.Parse(timeLayout, timestr)

		if err == nil {
			return &t, nil
		}
	}

	result := time.Now()

	return &result, errors.New("Cannot parse " + timestr)
}

func formatGPXTime(time *time.Time) string {
	if time == nil {
		return ""
	}
	if time.Year() <= 1 {
		// Invalid date:
		return ""
	}
	return time.Format(TIMELAYOUTS[0])
}

func ParseFile(fileName string) (*GPX, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return ParseBytes(bytes)
}

func ParseBytes(bytes []byte) (*GPX, error) {
	version, _ := guessGPXVersion(bytes)
	if version == "1.0" {
		g := gpx10.NewGpx()
		err := xml.Unmarshal(bytes, &g)
		if err != nil {
			return nil, err
		}

		return convertFromGpx10Models(g), nil
	} else if version == "1.1" {
		g := gpx11.NewGpx()
		err := xml.Unmarshal(bytes, &g)
		if err != nil {
			return nil, err
		}

		return convertFromGpx11Models(g), nil
	} else {
		return nil, errors.New("Invalid version:" + version)
	}
}

func ParseString(str string) (*GPX, error) {
	return ParseBytes([]byte(str))
}
