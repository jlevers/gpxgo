package gpx

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

const TIME_FORMAT = "2006-01-02T15:04:05Z"

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func assertEquals(t *testing.T, var1 interface{}, var2 interface{}) {
	if var1 != var2 {
		fmt.Println(var1, "not equals to", var2)
		t.Error("Not equals")
	}
}

func assertLinesEquals(t *testing.T, string1, string2 string) {
	lines1 := strings.Split(string1, "\n")
	lines2 := strings.Split(string2, "\n")
	for i := 0; i < min(len(lines1), len(lines2)); i++ {
		line1 := strings.Trim(lines1[i], " \n\r\t")
		line2 := strings.Trim(lines2[i], " \n\r\t")
		if line1 != line2 {
			t.Error("Line (#", i, ") different:", line1, "\nand:", line2)
			break
		}
	}
	if len(lines1) != len(lines2) {
		fmt.Println("String1:", string1)
		fmt.Println("String2:", string2)
		t.Error("String have a different number of lines", len(lines1), "and", len(lines2))
		return
	}
}

func assertNil(t *testing.T, var1 interface{}) {
	if var1 != nil {
		fmt.Println(var1)
		t.Error("nil!")
	}
}

func assertNotNil(t *testing.T, var1 interface{}) {
	if var1 == nil {
		fmt.Println(var1)
		t.Error("nil!")
	}
}

func TestParseGPXTimes(t *testing.T) {
	datetimes := []string{
		"2013-01-02T12:07:08Z",
		"2013-01-02 12:07:08Z",
		"2013-01-02T12:07:08",
		"2013-01-02T12:07:08.034Z",
		"2013-01-02 12:07:08.045Z",
		"2013-01-02T12:07:08.123",
	}
	for _, value := range datetimes {
		fmt.Println("datetime:", value)
		parsedTime, err := parseGPXTime(value)
		fmt.Println(parsedTime)
		assertNil(t, err)
		assertNotNil(t, parsedTime)
		assertEquals(t, parsedTime.Year(), 2013)
		assertEquals(t, parsedTime.Month(), time.January)
		assertEquals(t, parsedTime.Day(), 2)
		assertEquals(t, parsedTime.Hour(), 12)
		assertEquals(t, parsedTime.Minute(), 7)
		assertEquals(t, parsedTime.Second(), 8)
	}
}

func testDetectVersion(t *testing.T, fileName, expectedVersion string) {
	f, err := os.Open(fileName)
	fmt.Println("err=", err)
	contents, _ := ioutil.ReadAll(f)
	version, err := guessGPXVersion(contents)
	fmt.Println("Version=", version)
	if err != nil {
		t.Error("Can't detect 1.1 GPX, error=" + err.Error())
	}
	if version != expectedVersion {
		t.Error("Can't detect 1.1 GPX")
	}
}

func TestDetect11GPXVersion(t *testing.T) {
	testDetectVersion(t, "../test_files/gpx1.1_with_all_fields.gpx", "1.1")
}

func TestDetect10GPXVersion(t *testing.T) {
	testDetectVersion(t, "../test_files/gpx1.0_with_all_fields.gpx", "1.0")
}

func TestParseAndReparseGPX11(t *testing.T) {
	gpxDocuments := []*GPX{}

	{
		gpxDoc, err := ParseFile("../test_files/gpx1.1_with_all_fields.gpx")
		if err != nil || gpxDoc == nil {
			t.Error("Error parsing:" + err.Error())
		}
		gpxDocuments = append(gpxDocuments, gpxDoc)
		assertEquals(t, gpxDoc.Version, "1.1")

		// Test after reparsing
		xml, err := gpxDoc.ToXml(ToXmlParams{Version: "1.1", Indent: true})
		//fmt.Println(string(xml))
		if err != nil {
			t.Error("Error serializing to XML:" + err.Error())
		}
		gpxDoc2, err := ParseBytes(xml)
		assertEquals(t, gpxDoc2.Version, "1.1")
		if err != nil {
			t.Error("Error parsing XML:" + err.Error())
		}
		gpxDocuments = append(gpxDocuments, gpxDoc2)

		// TODO: ToString 1.0 and check again
	}

	for i := 1; i < len(gpxDocuments); i++ {
		fmt.Println("Testing gpx doc #", i)

		gpxDoc := gpxDocuments[i]

		executeSample11GpxAsserts(t, gpxDoc)

		// Tests after reparsing as 1.0
	}
}

func executeSample11GpxAsserts(t *testing.T, gpxDoc *GPX) {
	assertEquals(t, gpxDoc.Version, "1.1")
	assertEquals(t, gpxDoc.Creator, "...")
	assertEquals(t, gpxDoc.Name, "example name")
	assertEquals(t, gpxDoc.AuthorName, "author name")
	assertEquals(t, gpxDoc.AuthorEmail, "aaa@bbb.com")
	assertEquals(t, gpxDoc.Description, "example description")
	assertEquals(t, gpxDoc.AuthorLink, "http://link")
	assertEquals(t, gpxDoc.AuthorLinkText, "link text")
	assertEquals(t, gpxDoc.AuthorLinkType, "link type")
	assertEquals(t, gpxDoc.Copyright, "gpxauth")
	assertEquals(t, gpxDoc.CopyrightYear, "2013")
	assertEquals(t, gpxDoc.CopyrightLicense, "lic")
	assertEquals(t, gpxDoc.Link, "http://link2")
	assertEquals(t, gpxDoc.LinkText, "link text2")
	assertEquals(t, gpxDoc.LinkType, "link type2")
	assertEquals(t, gpxDoc.Time.Format(TIME_FORMAT), time.Date(2013, time.January, 01, 12, 0, 0, 0, time.UTC).Format(TIME_FORMAT))
	assertEquals(t, gpxDoc.Keywords, "example keywords")

	// Waypoints:
	assertEquals(t, len(gpxDoc.Waypoints), 2)
	assertEquals(t, gpxDoc.Waypoints[0].Latitude, 12.3)
	assertEquals(t, gpxDoc.Waypoints[0].Longitude, 45.6)
	assertEquals(t, gpxDoc.Waypoints[0].Elevation, 75.1)
	assertEquals(t, gpxDoc.Waypoints[0].Timestamp.Format(TIME_FORMAT), "2013-01-02T02:03:00Z")
	assertEquals(t, gpxDoc.Waypoints[0].MagneticVariation, "1.1")
	assertEquals(t, gpxDoc.Waypoints[0].GeoidHeight, "2.0")
	assertEquals(t, gpxDoc.Waypoints[0].Name, "example name")
	assertEquals(t, gpxDoc.Waypoints[0].Comment, "example cmt")
	assertEquals(t, gpxDoc.Waypoints[0].Description, "example desc")
	assertEquals(t, gpxDoc.Waypoints[0].Source, "example src")
	// TODO
	// Links       []GpxLink
	assertEquals(t, gpxDoc.Waypoints[0].Symbol, "example sym")
	assertEquals(t, gpxDoc.Waypoints[0].Type, "example type")
	assertEquals(t, gpxDoc.Waypoints[0].TypeOfGpsFix, "2d")
	assertEquals(t, gpxDoc.Waypoints[0].Satellites, 5)
	assertEquals(t, gpxDoc.Waypoints[0].HorizontalDilution, 6.0)
	assertEquals(t, gpxDoc.Waypoints[0].VerticalDilution, 7.0)
	assertEquals(t, gpxDoc.Waypoints[0].PositionalDilution, 8.0)
	assertEquals(t, gpxDoc.Waypoints[0].AgeOfDGpsData, 9.0)
	assertEquals(t, gpxDoc.Waypoints[0].DGpsId, 45)
	// TODO: Extensions

	assertEquals(t, gpxDoc.Waypoints[1].Latitude, 13.4)
	assertEquals(t, gpxDoc.Waypoints[1].Longitude, 46.7)

	// Routes:
	assertEquals(t, len(gpxDoc.Routes), 2)
	assertEquals(t, gpxDoc.Routes[0].Name, "example name")
	assertEquals(t, gpxDoc.Routes[0].Comment, "example cmt")
	assertEquals(t, gpxDoc.Routes[0].Description, "example desc")
	assertEquals(t, gpxDoc.Routes[0].Source, "example src")
	assertEquals(t, gpxDoc.Routes[0].Number, 7)
	assertEquals(t, gpxDoc.Routes[0].Type, "rte type")
	assertEquals(t, len(gpxDoc.Routes[0].Points), 3)
	// TODO: Link
	// TODO: Points
	assertEquals(t, gpxDoc.Routes[0].Points[0].Elevation, 75.1)
	fmt.Println("t=", gpxDoc.Routes[0].Points[0].Timestamp)
	assertEquals(t, gpxDoc.Routes[0].Points[0].Timestamp.Format(TIME_FORMAT), "2013-01-02T02:03:03Z")
	assertEquals(t, gpxDoc.Routes[0].Points[0].MagneticVariation, "1.2")
	assertEquals(t, gpxDoc.Routes[0].Points[0].GeoidHeight, "2.1")
	assertEquals(t, gpxDoc.Routes[0].Points[0].Name, "example name r")
	assertEquals(t, gpxDoc.Routes[0].Points[0].Comment, "example cmt r")
	assertEquals(t, gpxDoc.Routes[0].Points[0].Description, "example desc r")
	assertEquals(t, gpxDoc.Routes[0].Points[0].Source, "example src r")
	// TODO
	//assertEquals(t, gpxDoc.Routes[0].Points[0].Link, "http://linkrtept")
	//assertEquals(t, gpxDoc.Routes[0].Points[0].Text, "rtept link")
	assertEquals(t, gpxDoc.Routes[0].Points[0].Type, "example type r")
	assertEquals(t, gpxDoc.Routes[0].Points[0].Symbol, "example sym r")
	assertEquals(t, gpxDoc.Routes[0].Points[0].Type, "example type r")
	assertEquals(t, gpxDoc.Routes[0].Points[0].TypeOfGpsFix, "3d")
	assertEquals(t, gpxDoc.Routes[0].Points[0].Satellites, 6)
	assertEquals(t, gpxDoc.Routes[0].Points[0].HorizontalDilution, 7.0)
	assertEquals(t, gpxDoc.Routes[0].Points[0].VerticalDilution, 8.0)
	assertEquals(t, gpxDoc.Routes[0].Points[0].PositionalDilution, 9.0)
	assertEquals(t, gpxDoc.Routes[0].Points[0].AgeOfDGpsData, 10.0)
	assertEquals(t, gpxDoc.Routes[0].Points[0].DGpsId, 99)
	// TODO: Extensions

	assertEquals(t, gpxDoc.Routes[1].Name, "second route")
	assertEquals(t, gpxDoc.Routes[1].Description, "example desc 2")
	assertEquals(t, len(gpxDoc.Routes[1].Points), 2)

	// Tracks:
	assertEquals(t, len(gpxDoc.Tracks), 2)
	assertEquals(t, gpxDoc.Tracks[0].Name, "example name t")
	assertEquals(t, gpxDoc.Tracks[0].Comment, "example cmt t")
	assertEquals(t, gpxDoc.Tracks[0].Description, "example desc t")
	assertEquals(t, gpxDoc.Tracks[0].Source, "example src t")
	assertEquals(t, gpxDoc.Tracks[0].Number, 1)
	assertEquals(t, gpxDoc.Tracks[0].Type, "t")
	// TODO link

	assertEquals(t, len(gpxDoc.Tracks[0].Segments), 2)

	assertEquals(t, len(gpxDoc.Tracks[0].Segments[0].Points), 1)
	assertEquals(t, len(gpxDoc.Tracks[0].Segments[1].Points), 0)
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].Elevation, 11.1)
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].Timestamp.Format(TIME_FORMAT), "2013-01-01T12:00:04Z")
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].MagneticVariation, "12")
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].GeoidHeight, "13")
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].Name, "example name t")
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].Comment, "example cmt t")
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].Description, "example desc t")
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].Source, "example src t")
	// TODO link
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].Symbol, "example sym t")
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].Type, "example type t")
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].TypeOfGpsFix, "3d")
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].Satellites, 100)
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].HorizontalDilution, 101.0)
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].VerticalDilution, 102.0)
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].PositionalDilution, 103.0)
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].AgeOfDGpsData, 104.0)
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].DGpsId, 99)
	// TODO extensions
}

func TestParseAndReparseGPX10(t *testing.T) {
	gpxDocuments := []*GPX{}

	{
		gpxDoc, err := ParseFile("../test_files/gpx1.0_with_all_fields.gpx")
		if err != nil || gpxDoc == nil {
			t.Error("Error parsing:" + err.Error())
		}
		gpxDocuments = append(gpxDocuments, gpxDoc)
		assertEquals(t, gpxDoc.Version, "1.0")

		// Test after reparsing
		xml, err := gpxDoc.ToXml(ToXmlParams{Version: "1.0", Indent: true})
		//fmt.Println(string(xml))
		if err != nil {
			t.Error("Error serializing to XML:" + err.Error())
		}
		gpxDoc2, err := ParseBytes(xml)
		assertEquals(t, gpxDoc2.Version, "1.0")
		if err != nil {
			t.Error("Error parsing XML:" + err.Error())
		}
		gpxDocuments = append(gpxDocuments, gpxDoc2)

		// TODO: ToString 1.0 and check again
	}

	for i := 1; i < len(gpxDocuments); i++ {
		fmt.Println("Testing gpx doc #", i)

		gpxDoc := gpxDocuments[i]

		executeSample10GpxAsserts(t, gpxDoc)

		// Tests after reparsing as 1.0
	}
}

func executeSample10GpxAsserts(t *testing.T, gpxDoc *GPX) {
	assertEquals(t, gpxDoc.Version, "1.0")
	assertEquals(t, gpxDoc.Creator, "...")
	assertEquals(t, gpxDoc.Name, "example name")
	assertEquals(t, gpxDoc.AuthorName, "example author")
	assertEquals(t, gpxDoc.AuthorEmail, "example@email.com")
	assertEquals(t, gpxDoc.Description, "example description")
	assertEquals(t, gpxDoc.AuthorLink, "")
	assertEquals(t, gpxDoc.AuthorLinkText, "")
	assertEquals(t, gpxDoc.AuthorLinkType, "")
	assertEquals(t, gpxDoc.Copyright, "")
	assertEquals(t, gpxDoc.CopyrightYear, "")
	assertEquals(t, gpxDoc.CopyrightLicense, "")
	assertEquals(t, gpxDoc.Link, "http://example.url")
	assertEquals(t, gpxDoc.LinkText, "example urlname")
	assertEquals(t, gpxDoc.LinkType, "")
	assertEquals(t, gpxDoc.Time.Format(TIME_FORMAT), time.Date(2013, time.January, 01, 12, 0, 0, 0, time.UTC).Format(TIME_FORMAT))
	assertEquals(t, gpxDoc.Keywords, "example keywords")

	// TODO: Bounds (here and in 1.1)

	// Waypoints:
	assertEquals(t, len(gpxDoc.Waypoints), 2)
	assertEquals(t, gpxDoc.Waypoints[0].Latitude, 12.3)
	assertEquals(t, gpxDoc.Waypoints[0].Longitude, 45.6)
	assertEquals(t, gpxDoc.Waypoints[0].Elevation, 75.1)
	assertEquals(t, gpxDoc.Waypoints[0].Timestamp.Format(TIME_FORMAT), "2013-01-02T02:03:00Z")
	assertEquals(t, gpxDoc.Waypoints[0].MagneticVariation, "1.1")
	assertEquals(t, gpxDoc.Waypoints[0].GeoidHeight, "2.0")
	assertEquals(t, gpxDoc.Waypoints[0].Name, "example name")
	assertEquals(t, gpxDoc.Waypoints[0].Comment, "example cmt")
	assertEquals(t, gpxDoc.Waypoints[0].Description, "example desc")
	assertEquals(t, gpxDoc.Waypoints[0].Source, "example src")
	// TODO
	// Links       []GpxLink
	assertEquals(t, gpxDoc.Waypoints[0].Symbol, "example sym")
	assertEquals(t, gpxDoc.Waypoints[0].Type, "example type")
	assertEquals(t, gpxDoc.Waypoints[0].TypeOfGpsFix, "2d")
	assertEquals(t, gpxDoc.Waypoints[0].Satellites, 5)
	assertEquals(t, gpxDoc.Waypoints[0].HorizontalDilution, 6.0)
	assertEquals(t, gpxDoc.Waypoints[0].VerticalDilution, 7.0)
	assertEquals(t, gpxDoc.Waypoints[0].PositionalDilution, 8.0)
	assertEquals(t, gpxDoc.Waypoints[0].AgeOfDGpsData, 9.0)
	assertEquals(t, gpxDoc.Waypoints[0].DGpsId, 45)
	// TODO: Extensions

	assertEquals(t, gpxDoc.Waypoints[1].Latitude, 13.4)
	assertEquals(t, gpxDoc.Waypoints[1].Longitude, 46.7)

	// Routes:
	assertEquals(t, len(gpxDoc.Routes), 2)
	assertEquals(t, gpxDoc.Routes[0].Name, "example name")
	assertEquals(t, gpxDoc.Routes[0].Comment, "example cmt")
	assertEquals(t, gpxDoc.Routes[0].Description, "example desc")
	assertEquals(t, gpxDoc.Routes[0].Source, "example src")
	assertEquals(t, gpxDoc.Routes[0].Number, 7)
	assertEquals(t, gpxDoc.Routes[0].Type, "")
	assertEquals(t, len(gpxDoc.Routes[0].Points), 3)
	// TODO: Link
	// TODO: Points
	assertEquals(t, gpxDoc.Routes[0].Points[0].Elevation, 75.1)
	fmt.Println("t=", gpxDoc.Routes[0].Points[0].Timestamp)
	assertEquals(t, gpxDoc.Routes[0].Points[0].Timestamp.Format(TIME_FORMAT), "2013-01-02T02:03:03Z")
	assertEquals(t, gpxDoc.Routes[0].Points[0].MagneticVariation, "1.2")
	assertEquals(t, gpxDoc.Routes[0].Points[0].GeoidHeight, "2.1")
	assertEquals(t, gpxDoc.Routes[0].Points[0].Name, "example name r")
	assertEquals(t, gpxDoc.Routes[0].Points[0].Comment, "example cmt r")
	assertEquals(t, gpxDoc.Routes[0].Points[0].Description, "example desc r")
	assertEquals(t, gpxDoc.Routes[0].Points[0].Source, "example src r")
	// TODO link
	//assertEquals(t, gpxDoc.Routes[0].Points[0].Link, "http://linkrtept")
	//assertEquals(t, gpxDoc.Routes[0].Points[0].Text, "rtept link")
	assertEquals(t, gpxDoc.Routes[0].Points[0].Type, "example type r")
	assertEquals(t, gpxDoc.Routes[0].Points[0].Symbol, "example sym r")
	assertEquals(t, gpxDoc.Routes[0].Points[0].Type, "example type r")
	assertEquals(t, gpxDoc.Routes[0].Points[0].TypeOfGpsFix, "3d")
	assertEquals(t, gpxDoc.Routes[0].Points[0].Satellites, 6)
	assertEquals(t, gpxDoc.Routes[0].Points[0].HorizontalDilution, 7.0)
	assertEquals(t, gpxDoc.Routes[0].Points[0].VerticalDilution, 8.0)
	assertEquals(t, gpxDoc.Routes[0].Points[0].PositionalDilution, 9.0)
	assertEquals(t, gpxDoc.Routes[0].Points[0].AgeOfDGpsData, 10.0)
	assertEquals(t, gpxDoc.Routes[0].Points[0].DGpsId, 99)
	// TODO: Extensions

	assertEquals(t, gpxDoc.Routes[1].Name, "second route")
	assertEquals(t, gpxDoc.Routes[1].Description, "example desc 2")
	assertEquals(t, len(gpxDoc.Routes[1].Points), 2)

	// Tracks:
	assertEquals(t, len(gpxDoc.Tracks), 2)
	assertEquals(t, gpxDoc.Tracks[0].Name, "example name t")
	assertEquals(t, gpxDoc.Tracks[0].Comment, "example cmt t")
	assertEquals(t, gpxDoc.Tracks[0].Description, "example desc t")
	assertEquals(t, gpxDoc.Tracks[0].Source, "example src t")
	assertEquals(t, gpxDoc.Tracks[0].Number, 1)
	assertEquals(t, gpxDoc.Tracks[0].Type, "")
	// TODO link

	assertEquals(t, len(gpxDoc.Tracks[0].Segments), 2)

	assertEquals(t, len(gpxDoc.Tracks[0].Segments[0].Points), 1)
	assertEquals(t, len(gpxDoc.Tracks[0].Segments[1].Points), 0)
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].Elevation, 11.1)
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].Timestamp.Format(TIME_FORMAT), "2013-01-01T12:00:04Z")
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].MagneticVariation, "12")
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].GeoidHeight, "13")
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].Name, "example name t")
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].Comment, "example cmt t")
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].Description, "example desc t")
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].Source, "example src t")
	// TODO link
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].Symbol, "example sym t")
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].Type, "example type t")
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].TypeOfGpsFix, "3d")
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].Satellites, 100)
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].HorizontalDilution, 101.0)
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].VerticalDilution, 102.0)
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].PositionalDilution, 103.0)
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].AgeOfDGpsData, 104.0)
	assertEquals(t, gpxDoc.Tracks[0].Segments[0].Points[0].DGpsId, 99)
	// TODO extensions
}

func TestLength2DSeg(t *testing.T) {
	g, _ := ParseFile("../test_files/file.gpx")

	fmt.Println("tracks=", g.Tracks)
	fmt.Println("tracks=", len(g.Tracks))
	fmt.Println("segments=", len(g.Tracks[0].Segments))

	lengthA := g.Tracks[0].Segments[0].Length2D()
	lengthE := 56.77577732775905

	if lengthA != lengthE {
		t.Errorf("Length 2d expected: %f, actual %f", lengthE, lengthA)
	}
}

func TestLength3DSeg(t *testing.T) {
	g, _ := ParseFile("../test_files/file.gpx")
	lengthA := g.Tracks[0].Segments[0].Length3D()
	lengthE := 61.76815317436073

	if lengthA != lengthE {
		t.Errorf("Length 3d expected: %f, actual %f", lengthE, lengthA)
	}
}

func TestTimePoint(t *testing.T) {
	g, _ := ParseFile("../test_files/file.gpx")
	timeA := g.Tracks[0].Segments[0].Points[0].Timestamp
	//2012-03-17T12:46:19Z
	timeE := time.Date(2012, 3, 17, 12, 46, 19, 0, time.UTC)

	if timeA != timeE {
		t.Errorf("Time expected: %s, actual: %s", timeE.String(), timeA.String())
	}
}

func TestTimeBoundsSeg(t *testing.T) {
	g, _ := ParseFile("../test_files/file.gpx")
	timeBoundsA := g.Tracks[0].Segments[0].TimeBounds()

	startTime := time.Date(2012, 3, 17, 12, 46, 19, 0, time.UTC)
	endTime := time.Date(2012, 3, 17, 12, 47, 23, 0, time.UTC)
	timeBoundsE := TimeBounds{
		StartTime: startTime,
		EndTime:   endTime,
	}

	if !timeBoundsE.Equals(timeBoundsA) {
		t.Errorf("TimeBounds expected: %s, actual: %s", timeBoundsE.String(), timeBoundsA.String())
	}
}

func TestBoundsSeg(t *testing.T) {
	g, _ := ParseFile("../test_files/file.gpx")

	boundsA := g.Tracks[0].Segments[0].Bounds()
	boundsE := GpxBounds{
		MaxLat: 52.5117189623, MinLat: 52.5113534275,
		MaxLon: 13.4571944922, MinLon: 13.4567520116,
	}

	if !boundsE.Equals(boundsA) {
		t.Errorf("Bounds expected: %s, actual: %s", boundsE.String(), boundsA.String())
	}
}

func TestBoundsGpx(t *testing.T) {
	g, _ := ParseFile("../test_files/file.gpx")

	boundsA := g.Bounds()
	boundsE := GpxBounds{
		MaxLat: 52.5117189623, MinLat: 52.5113534275,
		MaxLon: 13.4571944922, MinLon: 13.4567520116,
	}

	if !boundsE.Equals(boundsA) {
		t.Errorf("Bounds expected: %s, actual: %s", boundsE.String(), boundsA.String())
	}
}

func TestSpeedSeg(t *testing.T) {
	g, _ := ParseFile("../test_files/file.gpx")
	speedA := g.Tracks[0].Segments[0].Speed(2)
	speedE := 1.5386074011963367

	if speedE != speedA {
		t.Errorf("Speed expected: %f, actual: %f", speedE, speedA)
	}
}

func TestSegmentDuration(t *testing.T) {
	g, _ := ParseFile("../test_files/file.gpx")
	durE := 64.0
	durA := g.Tracks[0].Segments[0].Duration()
	if durE != durA {
		t.Errorf("Duration expected: %f, actual: %f", durE, durA)
	}
}

func TestTrackDuration(t *testing.T) {
	g, _ := ParseFile("../test_files/file.gpx")
	durE := 64.0
	durA := g.Duration()
	if durE != durA {
		t.Errorf("Duration expected: %f, actual: %f", durE, durA)
	}
}

func TestMultiSegmentDuration(t *testing.T) {
	g, _ := ParseFile("../test_files/file.gpx")
	g.Tracks[0].AppendSegment(g.Tracks[0].Segments[0])
	durE := 64.0 * 2
	durA := g.Duration()
	if durE != durA {
		t.Errorf("Duration expected: %f, actual: %f", durE, durA)
	}
}

func TestMultiTrackDuration(t *testing.T) {
	g, _ := ParseFile("../test_files/file.gpx")

	g.Tracks[0].AppendSegment(g.Tracks[0].Segments[0])
	g.AppendTrack(g.Tracks[0])
	g.Tracks[0].AppendSegment(g.Tracks[0].Segments[0])

	durE := 384.0
	durA := g.Duration()
	if durE != durA {
		t.Errorf("Duration expected: %f, actual: %f", durE, durA)
	}
}

func TestUphillDownHillSeg(t *testing.T) {
	g, _ := ParseFile("../test_files/file.gpx")
	updoA := g.Tracks[0].Segments[0].UphillDownhill()
	updoE := UphillDownhill{
		Uphill:   5.863000000000007,
		Downhill: 1.5430000000000064}

	if !updoE.Equals(updoA) {
		t.Errorf("UphillDownhill expected: %+v, actual: %+v", updoE, updoA)
	}
}

func TestMovingData(t *testing.T) {
	g, _ := ParseFile("../test_files/file.gpx")
	movDataA := g.MovingData()
	movDataE := MovingData{
		MovingTime:      39.0,
		StoppedTime:     25.0,
		MovingDistance:  55.28705571308896,
		StoppedDistance: 6.481097461271765,
		MaxSpeed:        0.0,
	}

	if !movDataE.Equals(movDataA) {
		t.Errorf("Moving data expected: %+v, actual: %+v", movDataE, movDataA)
	}
}

func TestUphillDownhill(t *testing.T) {
	g, _ := ParseFile("../test_files/file.gpx")
	updoA := g.UphillDownhill()
	updoE := UphillDownhill{
		Uphill:   5.863000000000007,
		Downhill: 1.5430000000000064}

	if !updoE.Equals(updoA) {
		t.Errorf("UphillDownhill expected: %+v, actual: %+v", updoE, updoA)
	}
}

func TestToXml(t *testing.T) {
	g, _ := ParseFile("../test_files/file.gpx")
	xml, _ := g.ToXml(ToXmlParams{Version: "1.1", Indent: true})
	xmlA := string(xml)
	xmlE := `<?xml version="1.0" encoding="UTF-8"?>
<gpx xmlns="http://www.topografix.com/GPX/1/1" version="1.1" creator="eTrex 10">
	<metadata>
        <author></author>
		<link href="http://www.garmin.com">
			<text>Garmin International</text>
		</link>
		<time>2012-03-17T15:44:18Z</time>
	</metadata>
	<wpt lat="37.085751" lon="-121.17042">
		<ele>195.440933</ele>
		<time>2012-03-21T21:24:43Z</time>
		<name>001</name>
		<sym>Flag, Blue</sym>
	</wpt>
	<wpt lat="37.085751" lon="-121.17042">
		<ele>195.438324</ele>
		<time>2012-03-21T21:24:44Z</time>
		<name>002</name>
		<sym>Flag, Blue</sym>
	</wpt>
	<trk>
		<name>17-MRZ-12 16:44:12</name>
		<trkseg>
			<trkpt lat="52.5113534275" lon="13.4571944922">
				<ele>59.26</ele>
				<time>2012-03-17T12:46:19Z</time>
			</trkpt>
			<trkpt lat="52.5113568641" lon="13.4571697656">
				<ele>65.51</ele>
				<time>2012-03-17T12:46:44Z</time>
			</trkpt>
			<trkpt lat="52.511710329" lon="13.456941694">
				<ele>65.99</ele>
				<time>2012-03-17T12:47:01Z</time>
			</trkpt>
			<trkpt lat="52.5117189623" lon="13.4567520116">
				<ele>63.58</ele>
				<time>2012-03-17T12:47:23Z</time>
			</trkpt>
		</trkseg>
	</trk>
</gpx>`

	assertLinesEquals(t, xmlE, xmlA)
}

func TestNewXml(t *testing.T) {
	gpx := new(GPX)
	gpxTrack := new(GPXTrack)

	gpxSegment := new(GPXTrackSegment)
	gpxSegment.Points = append(gpxSegment.Points, &GPXPoint{Point: Point{Latitude: 2.1234, Longitude: 5.1234, Elevation: 1234}})
	gpxSegment.Points = append(gpxSegment.Points, &GPXPoint{Point: Point{Latitude: 2.1233, Longitude: 5.1235, Elevation: 1235}})
	gpxSegment.Points = append(gpxSegment.Points, &GPXPoint{Point: Point{Latitude: 2.1235, Longitude: 5.1236, Elevation: 1236}})

	gpxTrack.Segments = append(gpxTrack.Segments, gpxSegment)
	gpx.Tracks = append(gpx.Tracks, gpxTrack)

	xml, _ := gpx.ToXml(ToXmlParams{Version: "1.1", Indent: true})
	actualXml := string(xml)
	// TODO: xsi namespace:
	//expectedXml := `<gpx xmlns="http://www.topografix.com/GPX/1/1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd" version="1.1" creator="https://github.com/ptrv/go-gpx">
	expectedXml := `<?xml version="1.0" encoding="UTF-8"?>
<gpx xmlns="http://www.topografix.com/GPX/1/1" version="1.1" creator="https://github.com/ptrv/go-gpx">
	<metadata>
			<author></author>
	</metadata>
	<trk>
		<trkseg>
			<trkpt lat="2.1234" lon="5.1234">
				<ele>1234</ele>
			</trkpt>
			<trkpt lat="2.1233" lon="5.1235">
				<ele>1235</ele>
			</trkpt>
			<trkpt lat="2.1235" lon="5.1236">
				<ele>1236</ele>
			</trkpt>
		</trkseg>
	</trk>
</gpx>`

	assertLinesEquals(t, expectedXml, actualXml)
}

func TestInvalidXML(t *testing.T) {
	xml := "<gpx></gpx"
	gpx, err := ParseString(xml)
	if err == nil {
		t.Error("No error for invalid XML!")
	}
	if gpx != nil {
		t.Error("No gpx should be returned for invalid XMLs")
	}
}

func TestAddElevation(t *testing.T) {
	gpx := new(GPX)
	gpx.AppendTrack(new(GPXTrack))
	gpx.Tracks[0].AppendSegment(new(GPXTrackSegment))
	gpx.Tracks[0].Segments[0].AppendPoint(&GPXPoint{Point: Point{Latitude: 12, Longitude: 13, Elevation: 100}})
	gpx.Tracks[0].Segments[0].AppendPoint(&GPXPoint{Point: Point{Latitude: 12, Longitude: 13}})

	gpx.AddElevation(10)
	assertEquals(t, gpx.Tracks[0].Segments[0].Points[0].Elevation, 110.0)
	assertEquals(t, gpx.Tracks[0].Segments[0].Points[1].Elevation, 10.0) // TODO: this should be nil!

	gpx.AddElevation(-20)
	assertEquals(t, gpx.Tracks[0].Segments[0].Points[0].Elevation, 90.0)
	assertEquals(t, gpx.Tracks[0].Segments[0].Points[1].Elevation, -10.0) // TODO: this should be nil!
}

func TestRemoveElevation(t *testing.T) {
	g, _ := ParseFile("../test_files/file.gpx")

	g.RemoveElevation()

	xml, _ := g.ToXml(ToXmlParams{Indent: true})

	//fmt.Println(string(xml))

	if strings.Contains(string(xml), "<ele") {
		t.Error("Elevation still there!")
	}
}

func TestExecuteOnAllPoints(t *testing.T) {
	g, _ := ParseFile("../test_files/file.gpx")
	g.ExecuteOnAllPoints(func(*GPXPoint) {
	})
}