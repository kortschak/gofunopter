package display

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// Struct is a type for telling what display to display during optimization
type Struct struct {
	Value   interface{}
	Heading string
}

// Displayer is an interface for displaying values and headings
type Displayer interface {
	AddToDisplay([]*Struct) []*Struct
	Disp() bool
}

type Display struct {
	structs            []*Struct
	disp               bool
	headings           []string
	values             []string
	maxLengths         []int
	headingInterval    int
	valueInterval      time.Duration
	lastHeadingDisplay int // How many times have values been displayed since the last heading display
	lastValueDisplay   time.Time
}

func NewDisplay() *Display {
	return &Display{
		//structs: make([]*Struct, 0),
		disp:               true,
		headingInterval:    30,
		valueInterval:      500 * time.Millisecond,
		lastHeadingDisplay: math.MaxInt32 - 1, // High number so triggered on the first pass
		//lastValueDisplay:   Initialize to zero time so we print on the first iteration
	}
}

func (o *Display) IncreaseValueTime() {
	o.lastValueDisplay = time.Time{} //reset to zero time
}

// OptDisplay is so the optimizer implements Optimizer by
// embedding OptDisplay
func (o *Display) GetDisplay() *Display {
	return o
}

// Disp returns the toggle for the display of the whole optimizer
func (o *Display) Disp() bool {
	return o.disp
}

// SetDisp returns the toggle for the display of the whole optimizer
func (o *Display) SetDisp(b bool) {
	o.disp = b
}

// Note: This would ideally take in common.Common, but common imports display
// so it can't work the other way
// Note: userFun might not implement displayer, so need to check if it's nil

func (o *Display) DisplayProgress(displayers ...Displayer) {
	// Check that it's been long enough to display values
	if time.Since(o.lastValueDisplay) > o.valueInterval {
		// Collect all the values and headings to be displayed
		for _, displayer := range displayers {
			o.structs = o.structs[:0]
			if displayer.Disp() {
				o.structs = displayer.AddToDisplay(o.structs)
			}
			/*
				o.structs = o.structs[:0]
				if common.Disp() {
					o.structs = common.AddToDisplay(o.structs)
				}
				o.structs = optimizer.AddToDisplay(o.structs)
				if userFun != nil {
					o.structs = optimizer.AddToDisplay(o.structs)
				}
			*/
		}
		// Collect the print lengths of the values and headings and pad to match
		o.headings = o.headings[:0]
		o.values = o.values[:0]
		for _, str := range o.structs {
			var valueString string
			switch str.Value.(type) {
			case int:
				valueString = fmt.Sprintf("%d", str.Value)
			case float64:
				valueString = fmt.Sprintf("%e", str.Value)
			default:
				valueString = fmt.Sprintf("%v", str.Value)
			}
			if len(valueString) > len(str.Heading) {
				o.values = append(o.values, valueString)
				o.headings = append(o.headings, str.Heading+strings.Repeat(" ", len(valueString)-len(str.Heading)))
			} else {
				o.headings = append(o.headings, str.Heading)
				o.values = append(o.values, valueString+strings.Repeat(" ", len(str.Heading)-len(valueString)))
			}

		}

		// Print the values
		if o.lastHeadingDisplay > o.headingInterval {
			for _, val := range o.headings {
				fmt.Printf(val)
				fmt.Printf("\t")
			}
			fmt.Printf("\n")
			o.lastHeadingDisplay = 0
		}
		for _, val := range o.values {
			fmt.Printf(val)
			fmt.Printf("\t")
		}
		fmt.Printf("\n")
		o.lastHeadingDisplay++
		o.lastValueDisplay = time.Now()
	}
	////aoenustahoesnuthaonsteuhansoteihansoteuhanostu
}

/*
// TODO: Change appending values and headings to be one function

func AppendValues(values []interface{}, displayables ...Displayer) []interface{} {
	for _, displayable := range displayables {
		if displayable.Disp() {
			values = displayable.AppendValues(values)
		}
	}
	return values
}

func AppendHeadings(headings []string, displayables ...Displayer) []string {
	for _, displayable := range displayables {
		if displayable.Disp() {
			headings = displayable.AppendHeadings(headings)
		}
	}
	return headings
}

type Displayer interface {
	Disp() bool
	AppendHeadings([]string) []string
	AppendValues([]interface{}) []interface{}
	SetDisp(bool)
}

type GetDisplayStructer interface {
	Displayer
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
	return &DisplayStruct{
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

// Iterate the display. Checks to see if the columns or values should be displayed
func (d *DisplayStruct) Iterate() error {
	// TODO: Should this even be here?
	// TODO: Heading and value number error should give the headings and values
	if !d.DisplayOn {
		return nil // Display is off, don't do anything
	}
	headings := make([]string, 0)
	values := make([]interface{}, 0)
	d.headings = d.AppendHeadings(headings)
	d.values = d.AppendValues(values)
	if len(d.headings) != len(d.values) {
		fmt.Println("Here")
		fmt.Println("Headings", d.headings)
		fmt.Println("Values", d.values)
		return fmt.Errorf("Number of headings and values must match")
	}
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
*/
