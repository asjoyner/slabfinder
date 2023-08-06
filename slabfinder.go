package slabfinder

import (
	"time"
)

type Vendor int
type Finish int

const (
	StoneBasyx Vendor = 1

	UnknownFinish Finish = 0
	Polished      Finish = 1
	Leather       Finish = 2
	Honed         Finish = 3
)

// Slab describese one slab, which is in stock
type Slab struct {
	// TODO: add Material, if its ever other than Granite
	ID        int
	Price     int // in pennies
	Color     string
	Finish    Finish
	Thickness float64 // in CM
	Lot       string
	Bundle    string
	Width     float64 // inches
	Length    float64 // inches
	Count     int     // how many slabs are in this set
	Vendor    Vendor  // who has this slab for sale
	URL       string  // the detail page for the slab
	Photo     string  // the URL to a photo of the slab
	FirstSeen time.Time
	LastSeen  time.Time
}
