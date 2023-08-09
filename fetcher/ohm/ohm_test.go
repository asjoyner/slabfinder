package ohm

import (
	"os"
	"testing"

	"github.com/asjoyner/slabfinder"
	"github.com/google/go-cmp/cmp"
)

func TestParseJSON(t *testing.T) {
	type test struct {
		name     string
		input    string   // path to an html file
		SlabPage SlabPage // Metadata about the request
		want     []slabfinder.Slab
	}

	tests := []test{
		{
			name:     "Classic",
			input:    "testdata/titanium.charlotte.json",
			SlabPage: pages[0],
			want: []slabfinder.Slab{
				{
					Finish: slabfinder.Polished,
					Lot:    "6656",
					Bundle: "1497U",
					Width:  77.5,
					Length: 130,
					Count:  2,
					Vendor: slabfinder.Cosmos,
					URL:    "https://www.cosmosgranite.com/charlotte/granite/charlotte-293-titanium",
					Photo:  "https://cosmosgranite.nyc3.digitaloceanspaces.com/img/live_inventory/charlotte_charleston/LotImg_Titanium_6656_34969_1497U_A22.JPEG",
				},
				{
					Finish: slabfinder.Polished,
					Lot:    "8907",
					Bundle: "195320",
					Width:  77.5,
					Length: 131.5,
					Count:  4,
					Vendor: slabfinder.Cosmos,
					URL:    "https://www.cosmosgranite.com/charlotte/granite/charlotte-293-titanium",
					Photo:  "https://cosmosgranite.nyc3.digitaloceanspaces.com/img/live_inventory/charlotte_charleston/LotImg_Titanium_8907_36889_195320.JPEG",
				},
				{
					Finish: slabfinder.Polished,
					Lot:    "8907",
					Bundle: "195323",
					Width:  77,
					Length: 132,
					Count:  5,
					Vendor: slabfinder.Cosmos,
					URL:    "https://www.cosmosgranite.com/charlotte/granite/charlotte-293-titanium",
					Photo:  "https://cosmosgranite.nyc3.digitaloceanspaces.com/img/live_inventory/charlotte_charleston/LotImg_Titanium_8907_36889_195323.JPEG",
				},
				{
					Finish: slabfinder.Polished,
					Lot:    "8907",
					Bundle: "195523",
					Width:  75.5,
					Length: 121.5,
					Count:  5,
					Vendor: slabfinder.Cosmos,
					URL:    "https://www.cosmosgranite.com/charlotte/granite/charlotte-293-titanium",
					Photo:  "https://cosmosgranite.nyc3.digitaloceanspaces.com/img/live_inventory/charlotte_charleston/LotImg_Titanium_8907_36889_195523.JPEG",
				},
				{
					Finish: slabfinder.Polished,
					Lot:    "6135",
					Bundle: "349986",
					Width:  75,
					Length: 120.5,
					Count:  2,
					Vendor: slabfinder.Cosmos,
					URL:    "https://www.cosmosgranite.com/charlotte/granite/charlotte-293-titanium",
					Photo:  "https://cosmosgranite.nyc3.digitaloceanspaces.com/img/live_inventory/charlotte_charleston/LotImg_Titanium_6135_34021_349986.JPG",
				},
			},
		},
	}

	for _, tc := range tests {
		page, err := os.ReadFile(tc.input)
		if err != nil {
			t.Errorf("%s: %s", tc.name, err)
			continue
		}

		got, err := parseJSON(page, tc.SlabPage)
		if err != nil {
			t.Errorf("%s: parsing %s: %s", tc.name, tc.input, err)
			continue
		}

		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("%s:\n%s", tc.name, diff)
			continue
		}
	}
	return
}

/*
func TestParseFetch(t *testing.T) {
	t.Log(Fetch())
}
*/
