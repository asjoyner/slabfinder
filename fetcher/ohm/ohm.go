package ohm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/asjoyner/slabfinder"
)

// SlabPage defines the data necessary to fetch the Angular JSON data for types of slabs in a particular location
type SlabPage struct {
	Name     string
	FetchURL string
	Finish   slabfinder.Finish
}

var (
	pages = []SlabPage{
		{
			Name:     "AllSlabs",
			FetchURL: "https://ohm.stoneprofits.com/FetchDataWebV1.ashx?act=getItemGallery&InventoryGroupBy=IDTwo_&SearchbyItemIdentifiers=on&ShowFeatureProductOnTop=null&OnHold=null&OnSO=null&Intransit=null&showNotInStock=null&SearchbyFinish=on&SearchbySKU=on&Alphabet=",
			Finish:   slabfinder.Polished,
		},
		/*
			// This just fetches Copacabana slabs, which is a bit faster, but less flexible
			{
				Name:         "Copacabana",
				FetchURL:     "https://ohm.stoneprofits.com/FetchDataWebV1.ashx?act=getItemInventory&id=5181&InventoryGroupBy=IDTwo_&TrimmedUserID=4932186393528091&OnHold=null&OnSO=null&Intransit=null&SelectedLocation=&ShowLocationinGallery=on&LotPicturesRestrictToSIPL=False&ShowOnlyFullInventoryImages=on",
				Finish:       slabfinder.Polished,
			},
		*/
	}
)

// LinkURL is like https://inventory.ohmintl.com/CALCATTA-QUARTZITE-3CM-LEATHERED/4683/Location
// where CALCA.. is ItemName with dashes, and 4683 is ItemID
// PhotoBaseURL is always https://production123files.stoneprofits.com/Files/OHM

// Fetch consults all the OHM pages and returns the currently available slabs.
func Fetch() ([]slabfinder.Slab, error) {
	var slabs []slabfinder.Slab
	for _, page := range pages {
		resp, err := http.Get(page.FetchURL)
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

// SlabLotsJSONBody describes the body returned by the POST request
type SlabLotsJSONBody struct {
	SlabLot []SlabLot
}

// SlabLot describes one lot of slabs found at this page
type SlabLot struct {
	SELECTEDLocation string
	CategoryName     string
	ProductFormValue string
	ItemName         string
	ItemID           int
	FileName         string
	IDTwo            string
	Location         string
	LocationID       int
	CustomID         int
	FileID           string
	AverageLength    int
	AverageWidth     int
	AvailableQty     int
	UOM              string
	AvailableSlabs   int
	WebCartID        int
	Barcode          string
	Totalrows        string
}

// SlabTypesJSONBody describes the body returned by the POST request
type SlabTypesJSONBody struct {
	SlabType []SlabType
}

// SlabType describes one lot of slabs found at this page
type SlabType struct {
	Totalrows            string
	ItemID               int
	ItemName             string
	SKU                  string
	AlternateName        string
	DescriptiononWebsite string
	Origin               int    `json:",string"`
	Type                 string `json:"type"`
	TypeID               int    `json:",string"`
	Color                string
	NewArrival           string
	Filename             string
	CategoryName         string
	CategoryID           int
	SubCategory          string
	SubCategoryID        int
	LocationID           int
	Source               string
	PriceRange           string
	PriceRangeID         int
	GroupID              int `json:",string"`
	ThicknessID          int
	Thickness            int `json:",string"`
	ThicknessUOM         string
	ColorID              int `json:",string"`
	Finish               int
	OriginID             int `json:",string"`
	Kind                 string
	FeatureProduct       string
	IDTwo                int `json:",string"`
}

func parseJSON(body []byte, page SlabPage) ([]slabfinder.Slab, error) {
	var resp SlabTypesJSONBody
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal %s: %s", page.Name, err)
	}
	var slabs []slabfinder.Slab
	for _, s := range resp.SlabType {
		if s.LotBundlePicture == "" {
			continue
		}
		//photoURL, err := url.JoinPath(page.PhotoBaseURL, s.LotBundlePicture)
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
			//URL:    page.LinkURL,
			//Photo: photoURL,
		}
		slabs = append(slabs, slab)
	}

	return slabs, nil
}
