package stonebasyx

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"

	"github.com/asjoyner/slabfinder"
)

const (
	classic string = "https://www.stonebasyx.com/live-inventory/product-details/?selproductid=536"
	honed   string = "https://www.stonebasyx.com/live-inventory/product-details/?selproductid=690"
	leather string = "https://www.stonebasyx.com/live-inventory/product-details/?selproductid=712"
)

// Fetch consults all the StoneBasyx pages and returns the currently available slabs.
func Fetch() ([]slabfinder.Slab, error) {
	var slabs []slabfinder.Slab
	for _, url := range []string{classic, honed, leather} {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Errorf("fetching URL: %s", err)
			continue
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Errorf("reading Body: %s", err)
			continue
		}
		slabSubset, err := parseHTML(body, url)
		if err != nil {
			return nil, err
		}
		slabs = append(slabs, slabSubset...)
	}
	return slabs, nil
}

func parseHTML(page []byte, fetchURL string) ([]slabfinder.Slab, error) {
	var slabs []slabfinder.Slab
	var color string
	var finish slabfinder.Finish
	var thickness float64
	var foundContent bool
	var err error
	scanner := bufio.NewScanner(bytes.NewReader(page))
	for scanner.Scan() {
		line := scanner.Text()
		if !foundContent { // skip a lot of headers
			if strings.Contains(line, "<!-- write data here -->") {
				foundContent = true
			}
			continue
		}
		if color == "" && strings.Contains(line, ">Color: <") {
			color, err = parseKey(line)
			if err != nil {
				return nil, fmt.Errorf("color line invalid: %+v, %s", err, line)
			}
		}
		if finish == slabfinder.UnknownFinish && strings.Contains(line, ">Finish: <") {
			fs, err := parseKey(line)
			if err != nil {
				return nil, fmt.Errorf("Finish line invalid: %+v, %s", err, line)
			}
			switch fs {
			case "Polished":
				finish = slabfinder.Polished
			case "Honed":
				finish = slabfinder.Honed
			case "Leather":
				finish = slabfinder.Leather
			case "Leathered":
				finish = slabfinder.Leather
			}
		}
		if thickness == 0.0 && strings.Contains(line, ">Thickness: <") {
			ts, err := parseKey(line)
			if err != nil {
				return nil, fmt.Errorf("Thickness line invalid: %+v, %s", err, line)
			}
			thickness, err = strconv.ParseFloat(strings.Split(ts, " ")[0], 64)
			if err != nil {
				return nil, fmt.Errorf("Thickness line invalid: %+v, %s", err, line)
			}
		}

		if strings.Contains(line, "class=\"thumbpicsm2017\"") {
			slab := slabfinder.Slab{
				Vendor:    slabfinder.StoneBasyx,
				Color:     color,
				Finish:    finish,
				Thickness: thickness,
			}
			token := strings.Split(line, "\"")
			if len(token) < 4 {
				return nil, fmt.Errorf("photo line invalid: %+v", line)
			}
			urlWithSuffix, err := url.JoinPath(fetchURL, token[3])
			if err != nil {
				return nil, fmt.Errorf("crafting slab photo URL: %+v", line)
			}
			slab.Photo = strings.Split(urlWithSuffix, "?")[0]

			scanner.Scan() // throw away a line
			scanner.Scan() // fetch the second line
			line := scanner.Text()
			slab.Lot, err = parseKey(line)
			if err != nil {
				return nil, fmt.Errorf("Lot/Block line invalid: %+v, %s", err, line)
			}

			scanner.Scan() // throw away a line
			scanner.Scan() // fetch the second line
			line = scanner.Text()
			slab.Bundle, err = parseKey(line)
			if err != nil {
				return nil, fmt.Errorf("Bundle line invalid: %+v, %s", err, line)
			}

			scanner.Scan() // throw away a line
			scanner.Scan() // fetch the second line
			line = scanner.Text()
			size, err := parseKey(line)
			if err != nil {
				return nil, fmt.Errorf("Size line invalid: %+v, %s", err, line)
			}
			token = strings.Split(size, " ")
			if len(token) < 3 {
				return nil, fmt.Errorf("size too short: %+v", line)
			}
			length := token[0]
			width := token[2]
			if slab.Length, err = strconv.ParseFloat(strings.TrimSuffix(length, "L"), 64); err != nil {
				return nil, fmt.Errorf("length malformed: %q in line: %+v", length, line)
			}
			if slab.Width, err = strconv.ParseFloat(strings.TrimSuffix(width, "H"), 64); err != nil {
				return nil, fmt.Errorf("length malformed: %q in line: %+v", width, line)
			}

			scanner.Scan() // throw away a line
			scanner.Scan() // fetch the second line
			line = scanner.Text()
			count, err := parseKey(line)
			if err != nil {
				return nil, fmt.Errorf("In Stock line invalid: %+v, %s", err, line)
			}
			if slab.Count, err = strconv.Atoi(strings.TrimSuffix(count, " slabs")); err != nil {
				return nil, fmt.Errorf("slab count invalid: %+v, %s", err, line)
			}

			slabs = append(slabs, slab)
		}

	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanning the HTML page:", err)
	}

	return slabs, nil
}

// parseKey splits apart lines like this:
// <strong style="padding-left:20px;">Color: <a style="padding-left:5px;">Black, White</a></strong>^M
// and returns the string "Black, White"
func parseKey(line string) (string, error) {
	var value string
	rightSide := strings.Split(line, ">")
	if len(rightSide) < 3 {
		return "", fmt.Errorf("less than 3 > chars")
	}
	value = strings.Split(rightSide[2], "<")[0]
	return value, nil
}

func parseHTMLFancy(page []byte) ([]slabfinder.Slab, error) {
	doc, err := html.Parse(bytes.NewReader(page))
	if err != nil {
		return nil, fmt.Errorf("parsing login page: %s", err)
	}

	// recursively parse the HTML doc, node by node
	// find the form with id=form, and copy its hidden input fields
	//hiddenFields := make(map[string]string, 0)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "h4" {
				fmt.Println(n)
				/*
					var id string
					for _, a := range n.Attr {
						if a.Key == "id" {
							id = a.Val
						}
					}
					if id != "form" {
						return // skip forms with the wrong id
					}
				*/
			}
			/*
				if n.Data == "input" {
					attrMap := unpackAttrs(n.Attr)
					if attrMap["type"] == "hidden" {
						hiddenFields[attrMap["name"]] = attrMap["value"]
					}
				}
			*/
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return nil, nil
}
