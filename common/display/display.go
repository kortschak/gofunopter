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
}

func AddToDisplay(d []*Struct, disper ...Displayer) []*Struct {
	for _, val := range disper {
		d = val.AddToDisplay(d)
	}
	return d
}

type Display struct {
	structs            []*Struct
	Disp               bool
	headings           []string
	values             []string
	maxLengths         []int
	HeadingInterval    int
	ValueInterval      time.Duration
	lastHeadingDisplay int // How many times have values been displayed since the last heading display
	lastValueDisplay   time.Time
}

func (d *Display) SetSettings(disp *DisplaySettings) {
	d.HeadingInterval = disp.DisplayHeadingInterval
	d.ValueInterval = disp.DisplayValueInterval
	d.Disp = disp.Display
}

func NewDisplay() *Display {
	return &Display{
		lastHeadingDisplay: math.MaxInt32 - 1, // High number so triggered on the first pass
		//lastValueDisplay:   Initialize to zero time so we print on the first iteration
	}
}

type DisplaySettings struct {
	Display                bool          // Should the optimizer display at all
	DisplayHeadingInterval int           // How many value displays between each heading display
	DisplayValueInterval   time.Duration // How much time should pass between value displays
}

func NewDisplaySettings() *DisplaySettings {
	return &DisplaySettings{
		Display:                true,
		DisplayHeadingInterval: 30,
		DisplayValueInterval:   500 * time.Millisecond,
	}
}

func (o *Display) Reset() {
	o.lastHeadingDisplay = math.MaxInt32 - 1
	o.lastValueDisplay = time.Time{}
}

func (o *Display) IncreaseValueTime() {
	o.lastValueDisplay = time.Time{} //reset to zero time
}

// Note: This would ideally take in common.Common, but common imports display
// so it can't work the other way
// Note: userFun might not implement displayer, so need to check if it's nil

func (o *Display) DisplayProgress(displayers ...Displayer) {
	// If the display is off, don't do anything
	if o.Disp {
		// Check that it's been long enough to display values
		if time.Since(o.lastValueDisplay) > o.ValueInterval {
			o.structs = o.structs[:0]
			// Collect all the values and headings to be displayed
			for _, displayer := range displayers {
				o.structs = displayer.AddToDisplay(o.structs)
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
			if o.lastHeadingDisplay > o.HeadingInterval {
				fmt.Printf("\n")
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
	}
}
