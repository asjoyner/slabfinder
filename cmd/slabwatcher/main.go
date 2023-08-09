package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gtuk/discordwebhook"

	"github.com/asjoyner/slabfinder"
	"github.com/asjoyner/slabfinder/fetcher/cosmos"
	"github.com/asjoyner/slabfinder/fetcher/stonebasyx"
)

// TODO: Make these flags, default to OS config dir paths
const (
	slabFile = "/tmp/slabfinder.slabs.json"
	hookFile = "/tmp/slabfinder.webhook"
)

func main() {
	hb, err := os.ReadFile(hookFile)
	if err != nil {
		log.Printf("reading hookfile: %s", err)
	}
	hookURL := strings.TrimSpace(string(hb))

	slabs, err := loadSlabs(slabFile)
	if err != nil {
		log.Printf("reading known slabs: %s", err)
		os.Exit(1)
	}

	for {
		watch(slabs, hookURL)
	}
}

// SlabMap is a map from the Slab.ID() to Slab for easy lookup
type SlabMap map[uint64]slabfinder.Slab

// loadSlabs loads the known slabs from disk, returns a convenient map format
//
// TODO: Handle bootstrap condition by checking if the file exists, and
// accepting an arg to proceed if no file exists
func loadSlabs(slabFile string) (SlabMap, error) {
	input, err := os.ReadFile(slabFile)
	if err != nil {
		return nil, fmt.Errorf("reading known slabs: %s", err)
	}
	var ss []slabfinder.Slab
	if err := json.Unmarshal(input, &ss); err != nil {
		return nil, fmt.Errorf("parsing known slabs: %s", err)
	}
	slabs := make(SlabMap)
	for _, slab := range ss {
		slabs[slab.ID()] = slab // recompute the ID each time, so changing it is less cumbersome
	}
	return slabs, nil
}

func watch(slabs SlabMap, hookURL string) {
	// Fetch the latest slabs
	// TODO: generalize this to iterate over all the vendors
	thisRunTimestamp := time.Now()
	var ns []slabfinder.Slab
	s, err := stonebasyx.Fetch()
	if err != nil {
		log.Println(err)
	}
	ns = append(ns, s...)
	s, err = cosmos.Fetch()
	if err != nil {
		log.Println(err)
	}
	ns = append(ns, s...)

	// include new slabs in the known slabs, update timestamps
	for _, slab := range ns {
		if oldSlab, ok := slabs[slab.ID()]; ok {
			slab.FirstSeen = oldSlab.FirstSeen
		} else {
			slab.FirstSeen = thisRunTimestamp
		}
		slab.LastSeen = thisRunTimestamp
		slabs[slab.ID()] = slab
	}

	// snapshot slabs to disk
	var wss []slabfinder.Slab
	for _, s := range slabs {
		wss = append(wss, s)
	}
	output, err := json.MarshalIndent(wss, "", "	")
	if err != nil {
		log.Printf("could not marshal slabs: %s", err)
		os.Exit(3)
	}
	if err := ioutil.WriteFile(slabFile, output, 0644); err != nil {
		log.Printf("writing slabs: %s", err)
		os.Exit(4)
	}

	// filter slabs by criteria
	var ourSlabs []slabfinder.Slab
	for _, slab := range slabs {
		if slab.Length < 132 {
			continue
		}
		if slab.FirstSeen == slab.LastSeen {
			fmt.Printf("Interesting new slab: %s\n", slab.String())
		}
		ourSlabs = append(ourSlabs, slab)
	}

	// TODO: write HTML page of known interesting slabs?

	// send Discord notification of new interesting slabs
	if hookURL != "" {
		discordUsername := "SlabFinder"
		for _, slab := range ourSlabs {
			if slab.FirstSeen != slab.LastSeen {
				continue
			}
			content := slab.String()
			image := discordwebhook.Image{Url: &slab.Photo}
			embed := discordwebhook.Embed{Image: &image}
			msg := discordwebhook.Message{
				Username: &discordUsername,
				Content:  &content,
				Embeds:   &[]discordwebhook.Embed{embed},
			}
			if err := discordwebhook.SendMessage(hookURL, msg); err != nil {
				fmt.Println(err)
			}
		}
	}
	time.Sleep(15 * time.Minute)
}
