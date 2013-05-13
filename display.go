package gofunopter

import (
	"fmt"
	"math"
	"time"
)

func AppendValues(values []interface{}, displayables ...Displayable) []interface{} {
	for _, displayable := range displayables {
		if displayable.Disp() {
			values = displayable.AppendValues(values)
		}
	}
	return values
}

func AppendHeadings(headings []string, displayables ...Displayable) []string {
	for _, displayable := range displayables {
		headings = displayable.AppendHeadings(headings)
	}
	return headings
}

type Displayer interface {
	Disp() bool
	SetDisp(bool)
}

type BasicDisplayer struct {
	disp bool
}

func (b *BasicDisplayer) Disp() bool {
	return b.disp
}

func (b *BasicDisplayer) SetDisp(val bool) {
	b.disp = val
}

func NewDisplay(val bool) *BasicDisplayer {
	return &BasicDisplayer{disp: val}
}

// Something which can display values
type Displayable interface {
	AppendHeadings([]string) []string
	AppendValues([]interface{}) []interface{}
	Disp() bool
	//Display() *Display
}

type GetDisplayStructer interface {
	Displayable
	GetDisplayStruct() *DisplayStruct
}

func SetDisplayMethods(h GetDisplayStructer) {
	d := h.GetDisplayStruct()
	d.AppendHeadings = h.AppendHeadings
	d.AppendValues = h.AppendValues
}

// TODO: Somehow turn this into a writer

// Controls the display settings for an optimizer
// On is a bool setting if values are displayed at all
// Time interval is a setting for at least how many seconds should elapse between value displays
// HeadingInterval sets how many value displays happen between reprinting the columns
// D is a displayer which (usually) should be set by the optimization algorithm, though customizations are possible
type DisplayStruct struct {
	DisplayOn       bool
	lastValueTime   time.Time     // When was the last real time the values were displayed
	nValueDisplays  int           // How many displays have there been since the headings were output
	TimeInterval    time.Duration // How many seconds should elapse between displays
	HeadingInterval int           // How many value outputs between displays
	headings        []string
	values          []interface{}
	//headinglengths  []int // For aligning the columns
	//valuelengths    []int
	AppendHeadings func([]string) []string
	AppendValues   func([]interface{}) []interface{}
	//D               Displayer
	//TODO: Add in something about the column widths
}

// Returns the default settings for the display parameters
func DefaultDisplayStruct() *DisplayStruct {
	// Defaults are for forcing display on the first iteration
	return &Display{
		DisplayOn:       true,
		TimeInterval:    700 * time.Millisecond,
		HeadingInterval: 30,
		lastValueTime:   time.Now().Add(math.MaxUint16 - 1),
		nValueDisplays:  math.MaxUint16 - 1,
	}
}

func (d *DisplayStruct) GetDisplayStruct() *DisplayStruct {
	return d
}

/*
func (d *Display) SetHeadings(s []string) {
	d.headings = s
}

func (d *Display) SetValues(v []interface{}) {
	d.values = v
}
*/

// Iterate the display. Checks to see if the columns or values should be displayed
func (d *DisplayStruct) Iterate() error {
	if !d.DisplayOn {
		return nil // Display is off, don't do anything
	}
	headings := make([]string, 0)
	values := make([]interface{}, 0)
	d.headings = d.AppendHeadings(headings)
	d.values = d.AppendValues(values)
	if len(d.headings) != len(d.values) {
		return fmt.Errorf("Number of headings and values must match")
	}
	//if cap(d.headinglengths) < len(d.headings) {
	//	d.headinglengths = make([]int, len(d.headings))
	//}
	//if cap(d.valuelengths) < len(d.values) {

	//}
	// Find how many characters it would take to display each of the things

	// First, check if the headings need to be set
	if d.nValueDisplays >= d.HeadingInterval {
		d.Headings()
		d.nValueDisplays = 0
	}
	// Then, check if the values need to be set
	if time.Since(d.lastValueTime) > d.TimeInterval {
		d.Values()
		d.nValueDisplays++
		d.lastValueTime = time.Now()
	}
	return nil
}

// Display the headings returned by the displayer
func (d *DisplayStruct) Headings() {
	fmt.Print("\n")
	for _, val := range d.headings {
		fmt.Print(val)
		fmt.Print("\t")
	}
	fmt.Print("\n")
}

func (d *DisplayStruct) Values() {
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
