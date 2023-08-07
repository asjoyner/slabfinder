package slabfinder

import (
	"fmt"
	"time"

	"github.com/cespare/xxhash"
)

type Vendor int
type Finish int

const (
	UnknownVendor Vendor = 0
	StoneBasyx    Vendor = 1
	Cosmos        Vendor = 2

	UnknownFinish Finish = 0
	Polished      Finish = 1
	Leather       Finish = 2
	Honed         Finish = 3
)

// Slab describese one slab, which is in stock
type Slab struct {
	// TODO: add Material, if its ever other than Granite
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

func (s *Slab) ID() uint64 {
	return xxhash.Sum64([]byte(fmt.Sprintf("%s%s%vd%s%s%s%s", s.Vendor, s.Finish, s.Thickness, s.Color, s.Lot, s.Bundle, s.Photo)))
}

func (s *Slab) String() string {
	return fmt.Sprintf("Length: %v, Count: %d, Lot: %s, Bundle: %s, Finish: %s, Vendor: %s, URL: %s", s.Length, s.Count, s.Lot, s.Bundle, s.Finish, s.Vendor, s.URL)
}

func (v Vendor) String() string {
	switch v {
	case StoneBasyx:
		return "StoneBasyx"
	case Cosmos:
		return "Cosmos"
	}
	return "UnknownVendor"
}

func (f Finish) String() string {
	switch f {
	case Polished:
		return "Polished"
	case Leather:
		return "Leather"
	case Honed:
		return "Honed"
	}
	return "UnknownPolish"
}
