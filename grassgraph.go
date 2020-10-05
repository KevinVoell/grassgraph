package grassgraph

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"

	"github.com/labstack/gommon/log"
)

/*
GetGrassGraph Gets the grass graph from github and converts it into a PNG image
*/
func GetGrassGraph(username string) ([]byte, error) {

	body, err := getSvgFromGithub(username)
	if err != nil {
		return nil, err
	}

	cleanedSvg, err := extractSVGAndFixup(body)
	if err != nil {
		return nil, err
	}

	pngBytes, err := convertSvgToPng(cleanedSvg)
	if err != nil {
		return nil, err
	}

	return pngBytes, nil
}

/*
getSvgFromGithub gets the SVG markup from github for the grass graph.
*/
func getSvgFromGithub(username string) (string, error) {
	requestURI := fmt.Sprintf("https://www.github.com/%s", username)

	response, err := http.Get(requestURI)
	if err != nil {
		log.Error("Failed to get data from github")
		return "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("Failed to read response body")
		return "", err
	}

	return string(body), nil
}

func extractSVGAndFixup(body string) (string, error) {
	repexp := regexp.MustCompile(`^[\s\S]+<svg.+class="js-calendar-graph-svg">`)
	repcnd := `<svg xmlns="http://www.w3.org/2000/svg" width="870" height="155" class="js-calendar-graph-svg">
		<rect x="0" y="0" width="828" height="128" fill="white" stroke="none"/>`
	graphData := repexp.ReplaceAllString(body, repcnd)

	repexp = regexp.MustCompile(`<text text-anchor="start" class="wday" dx="-10" dy="8" style="display: none;">Sun</text>`)
	repcnd = ``
	graphData = repexp.ReplaceAllString(graphData, repcnd)

	repexp = regexp.MustCompile(`<text text-anchor="start" class="wday" dx="-10" dy="32" style="display: none;">Tue</text>`)
	repcnd = ``
	graphData = repexp.ReplaceAllString(graphData, repcnd)

	repexp = regexp.MustCompile(`<text text-anchor="start" class="wday" dx="-10" dy="57" style="display: none;">Thu</text>`)
	repcnd = ``
	graphData = repexp.ReplaceAllString(graphData, repcnd)

	repexp = regexp.MustCompile(`<text text-anchor="start" class="wday" dx="-10" dy="81" style="display: none;">Sat</text>`)
	repcnd = ``
	graphData = repexp.ReplaceAllString(graphData, repcnd)

	repexp = regexp.MustCompile(`dy="81" style="display: none;">Sat<\/text>[\s\S]+<\/g>[\s\S]+<\/svg>[.\s\S]+\z`)
	repcnd = `dy="81" style="display: none;">Sat</text>
		<text x="675" y="125">Less</text>
		<g transform="translate(709,15)">
			<rect width="11" height="11" x="0" y="99" fill="#ebedf0"/>
		</g>
		<g transform="translate(724,15)">
			<rect width="11" height="11" y="99" fill="#9be9a8"/>
		</g>
		<g transform="translate(739,15)">
			<rect width="11" height="11" y="99" fill="#40c463"/>
		</g>
		<g transform="translate(754,15)">
			<rect width="11" height="11" y="99" fill="#30a14e"/>
		</g>
		<g transform="translate(769,15)">
			<rect width="11" height="11" y="99" fill="#216e39"/>
		</g>
		<text x="788" y="125">More</text>
	</g>
	</svg>`
	graphData = repexp.ReplaceAllString(graphData, repcnd)

	repexp = regexp.MustCompile(`fill="#ebedf0"`)
	//repcnd = `#d3d5d8`
	repcnd = `style="fill:white;stroke:#bcbdc0;stroke-width:1"`
	graphData = repexp.ReplaceAllString(graphData, repcnd)

	return graphData, nil
}

func convertSvgToPng(svgData string) ([]byte, error) {

	svgFilename := fmt.Sprintf("%s/github.svg", os.TempDir())

	err := ioutil.WriteFile(svgFilename, []byte(svgData), 0644)
	if err != nil {
		return nil, err
	}
	defer deleteFile(svgFilename)

	pngFilename := fmt.Sprintf("%s/github.png", os.TempDir())
	defer deleteFile(pngFilename)

	cmd := exec.Command("convert", "-geometry", "870x155", "-rotate", "0", svgFilename, pngFilename)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadFile(pngFilename)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func deleteFile(filename string) {
	_, err := os.Stat(filename)
	if err != nil {
		return
	}

	err = os.Remove(filename)
	if err != nil {
		return
	}
}
