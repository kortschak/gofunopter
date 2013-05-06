package gofunopter

import (
	"fmt"
	"time"
)

// Controls the display settings for an optimizer
// On is a bool setting if values are displayed at all
// Time interval is a setting for at least how many seconds should elapse between value displays
// HeadingInterval sets how many value displays happen between reprinting the columns
// D is a displayer which (usually) should be set by the optimization algorithm, though customizations are possible
type Display struct {
	DisplayOn       bool
	lastValueTime   time.Time // When was the last real time the values were displayed
	nValueDisplays  int       // How many displays have there been since the headings were output
	TimeInterval    float64   // How many seconds should elapse between displays
	HeadingInterval int       // How many value outputs between displays
	headings        []string
	values          []interface{}
	//D               Displayer
	//TODO: Add in something about the column widths
}

// Returns the default settings for the display parameters
func DefaultDisplay() *Display {
	return &Display{DisplayOn: true, TimeInterval: 0.7, HeadingInterval: 30}
}

func (d *Display) Display() *Display {
	return d
}

func (d *Display) SetHeadings(s []string) {
	d.headings = s
}

func (d *Display) SetValues(v []interface{}) {
	d.values = v
}

// Iterate the display. Checks to see if the columns or values should be displayed
func (d *Display) Iterate() {
	if !d.DisplayOn {
		return // Display is off, don't do anything
	}
	// First, check if the headings need to be set
	if d.nValueDisplays >= d.HeadingInterval {
		d.Headings()
		d.nValueDisplays = 0
	}
	// Then, check if the values need to be set
	if time.Since(d.lastValueTime).Seconds() > d.TimeInterval {
		d.Values()
		d.nValueDisplays++
		d.lastValueTime = time.Now()
	}
}

// Display the headings returned by the displayer
func (d *Display) Headings() {
	fmt.Print("\n")
	for _, val := range d.headings {
		fmt.Print(val)
		fmt.Print("\t")
	}
	fmt.Print("\n")
}

func (d *Display) Values() {
	for _, val := range d.values {
		switch val.(type) {
		case int:
			fmt.Printf("%d", val)
		case float64:
			fmt.Printf("%e", val)
		default:
			fmt.Printf("%v", val)
		}
		fmt.Printf("\t")
	}
	fmt.Print("\n")
}
