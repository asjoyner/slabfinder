package cosmos

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/asjoyner/slabfinder"
)

// SlabPage defines the data necessary to fetch the Angular JSON data for types of slabs in a particular location
type SlabPage struct {
	Name         string
	FetchURL     string
	PostData     string
	LinkURL      string
	PhotoBaseURL string
	Finish       slabfinder.Finish
}

var (
	pages = []SlabPage{
		{
			Name:         "Titanium",
			FetchURL:     "https://www.cosmosgranite.com/getProductDetail",
			PostData:     "name=Titanium&location=charlotte&id=20488&pro_link=https%3A%2F%2Fwww.cosmosgranite.com%2Fcharlotte%2Fgranite%2Fcharlotte-293-titanium",
			LinkURL:      "https://www.cosmosgranite.com/charlotte/granite/charlotte-293-titanium",
			PhotoBaseURL: "https://cosmosgranite.nyc3.digitaloceanspaces.com/img/live_inventory/charlotte_charleston/",
			Finish:       slabfinder.Polished,
		},
		{
			Name:         "Titanium Leathered",
			FetchURL:     "https://www.cosmosgranite.com/getProductDetail",
			PostData:     "urls=http%3A%2F%2Fapi.vividgranite.com%2Fservices.asmx&name=TITANIUM+LEATHER&lot=&bundle=&location=charlotte&id=30427&pro_link=https%3A%2F%2Fwww.cosmosgranite.com%2Fcharlotte%2Fgranite%2Fcharlotte-1311-titanium-leather",
			LinkURL:      "https://www.cosmosgranite.com/charlotte/granite/charlotte-1311-titanium-leather",
			PhotoBaseURL: "https://cosmosgranite.nyc3.digitaloceanspaces.com/img/live_inventory/charlotte_charleston/",
			Finish:       slabfinder.Polished,
		},
	}
)

// Fetch consults all the Cosmos pages and returns the currently available slabs.
func Fetch() ([]slabfinder.Slab, error) {
	var slabs []slabfinder.Slab
	for _, page := range pages {
		req, err := http.NewRequest("POST", page.FetchURL, strings.NewReader(page.PostData))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("x-requested-with", "XMLHttpRequest")
		//o, _ := httputil.DumpRequestOut(req, true)
		//fmt.Println(string(o))
		resp, err := http.DefaultClient.Do(req)
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
		slabSubset, err := parseJSON(body, page)
		if err != nil {
			return nil, err
		}
		slabs = append(slabs, slabSubset...)
	}
	return slabs, nil
}

// JSONBody describes the body returned by the POST request
type JSONBody struct {
	Msg    string      `json:"username"`
	Status int         `json:"status"`
	Slabs  []CosmoSlab `json:"api_data"`
}

// CosmoSlab describes one lot of slabs found at this page
type CosmoSlab struct {
	AvgSlabLength         float64 `json:",string"`
	AvgSlabWidth          float64 `json:",string"`
	ProductID             string
	ProductName           string
	AvailableSlabs        int     `json:",string"`
	AvailableQuantity     float64 `json:",string"`
	AvaialbleLocationName string
	ProductStatus         string
	LotNumber             string
	LotBundlePicture      string
	BundleNumber          string
	PSDUnique1            string `json:"PSD_Unique1"`
}

func parseJSON(body []byte, page SlabPage) ([]slabfinder.Slab, error) {
	var resp JSONBody
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal %s: %s", page.Name, err)
	}
	var slabs []slabfinder.Slab
	for _, s := range resp.Slabs {
		if s.LotBundlePicture == "" {
			continue
		}
		photoURL, err := url.JoinPath(page.PhotoBaseURL, s.LotBundlePicture)
		if err != nil {
			fmt.Errorf("invalid photo URL: %s", err)
		}
		slab := slabfinder.Slab{
			Finish: page.Finish,
			Lot:    s.LotNumber,
			Bundle: s.BundleNumber,
			Width:  s.AvgSlabWidth,
			Length: s.AvgSlabLength,
			Count:  s.AvailableSlabs,
			Vendor: slabfinder.Cosmos,
			URL:    page.LinkURL,
			Photo:  photoURL,
		}
		slabs = append(slabs, slab)
	}

	return slabs, nil
}
