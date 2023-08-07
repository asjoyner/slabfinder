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
	"github.com/asjoyner/slabfinder/fetcher/stonebasyx"
)

const (
	slabFile = "/tmp/slabfinder.slabs.json"
	hookFile = "/tmp/slabfinder.webhook"
)

func main() {
	hookURL, err := os.ReadFile(hookFile)
	if err != nil {
		log.Printf("reading hookfile: %s", err)
	}

	// Load known slabs
	input, err := os.ReadFile(slabFile)
	if err != nil {
		log.Printf("reading known slabs: %s", err)
		os.Exit(1)
	}
	var ss []slabfinder.Slab
	if err := json.Unmarshal(input, &ss); err != nil {
		log.Printf("parsing known slabs: %s", err)
		os.Exit(2)
	}
	slabs := make(map[uint64]slabfinder.Slab)
	for _, slab := range ss {
		slabs[slab.ID()] = slab // recompute the ID each time, so changing it is less cumbersome
	}

	// Fetch the latest slabs
	// TODO: generalize this to iterate over all the fetchers
	thisRunTimestamp := time.Now()
	ns, err := stonebasyx.Fetch()
	if err != nil {
		log.Println(err)
	}

	// include new entries in the known slabs, update timestamps
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

	// write HTML page of known interesting slabs
	// send notification of new interesting slabs
	if hookURL != nil {
		discordUsername := "SlabFinder"
		hk := strings.TrimSpace(string(hookURL))
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
			if err := discordwebhook.SendMessage(hk, msg); err != nil {
				fmt.Println(err)
			}
		}
	}
}
