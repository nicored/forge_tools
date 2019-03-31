package properties

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

// Constants representing the type of data in each column of an attribute row.
const (
	aName             = iota // 0
	aCategory                // 1
	aType                    // 2
	aUnit                    // 3
	aDesc                    // 4
	aDisplayName             // 5
	aFlags                   // 6
	aDisplayPrecision        // 7
)

type Attr struct {
	Idx              int
	Name             string
	Category         string
	Type             *int
	Unit             *string
	Description      *string
	DisplayName      *string
	Flags            bool
	DisplayPrecision bool
}

type Attrs []*Attr

// readAttrs returns a slice of Attr representing
// the attribute properties as defined in objects_attrs.json
func readAttrs(path string) Attrs {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	attrs, err := ioutil.ReadAll(f)

	var m = &[]interface{}{}
	err = json.Unmarshal(attrs, m)
	if err != nil {
		panic(err)
	}

	var attrList Attrs

	for i, r := range *m {
		if _, ok := r.(float64); ok {
			continue
		}

		a := parseAttr(r)
		a.Idx = i
		attrList = append(attrList, &a)
	}

	return attrList
}

// parseAttr creates and populates an Attr struct
// with the data provided in the given row r
func parseAttr(r interface{}) Attr {
	row := r.([]interface{})

	a := Attr{
		Name:             row[aName].(string),
		Type:             nilInt(row[aType]),
		Unit:             nilString(row[aUnit]),
		Description:      nilString(row[aDesc]),
		DisplayName:      nilString(row[aDisplayName]),
		Flags:            rBool(row[aFlags]),
		DisplayPrecision: false,
	}

	// We don't want empty or nil categories, so we change them
	// to 'General' and 'Attribute'
	c := nilString(row[aCategory])
	if c == nil {
		cd := "General"
		c = &cd
	} else if *c == "" {
		cd := "Attribute"
		c = &cd
	}
	a.Category = *c

	if len(row) > 7 {
		a.DisplayPrecision = rBool(row[aDisplayPrecision])
	}

	return a
}

// List of integers from the objects_attrs.json file
type Offs []int

// Returns the max ID from the Offs data set.
// As ID's start at 0, the max ID is the length of Offs array - 1
func (o Offs) IdMax() int {
	return len(o) - 1
}

// readOffs returns a list of integers from the objects_offs.json file
func readOffs(path string) Offs {
	return readIntSlice(path)
}

// Avs is the
type Avs []int

func readAvs(path string) Avs {
	return readIntSlice(path)
}

type Vals []interface{}

func readVals(path string) Vals {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	attrs, err := ioutil.ReadAll(f)

	var m = Vals{}
	err = json.Unmarshal(attrs, &m)
	if err != nil {
		panic(err)
	}

	return m
}

type Ids []string

func readIds(path string) Ids {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	attrs, err := ioutil.ReadAll(f)

	var m []interface{}
	err = json.Unmarshal(attrs, &m)
	if err != nil {
		panic(err)
	}

	ids := Ids{}
	for _, id := range m {
		// ids are hexadecimal numbers in string
		// the json file does not show consistency and does not usually
		// quote the value at the first index '0', which is therefore read as a float64
		if vid, ok := id.(string); ok {
			ids = append(ids, vid)
		} else if vid, ok := id.(float64); ok {
			ivid := int(vid)
			ids = append(ids, strconv.Itoa(ivid))
		}
	}

	return ids
}

// readIntSlice reads a json file as a slice of integers
func readIntSlice(path string) []int {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	offs, err := ioutil.ReadAll(f)

	var m []int
	err = json.Unmarshal(offs, &m)
	if err != nil {
		panic(err)
	}

	return m
}

// nilString returns a string pointer. If the interface value
// is nil, it returns a nil string, else returns a pointer to the
// string value
func nilString(v interface{}) *string {
	if vs, ok := v.(string); ok {
		return &vs
	} else if v == nil {
		return nil
	}

	panic("not a string")
}

// nilInt returns a integer pointer. If the interface value
// is nil, it returns a nil integer, else returns a pointer to the
// integer value. If the values were decoded from a json file, numbers
// will be of type float64, which we convert to int
func nilInt(v interface{}) *int {
	if vf, ok := v.(float64); ok {
		vi := int(vf)
		return &vi
	} else if vi, ok := v.(int); ok {
		return &vi
	}

	panic("not a number")
}

// rBool returns a boolean.
// if the interface value is a number, any value greater than 0 will return true, else false
// if the interface value is a string, a non-empty string will return true, else false
// if the interface value is null, the function will return false
func rBool(v interface{}) bool {
	if v == nil {
		return false
	} else if vf, ok := v.(float64); ok {
		return vf > 0
	} else if vi, ok := v.(int); ok {
		return vi > 0
	} else if vs, ok := v.(string); ok {
		return vs != ""
	} else if vb, ok := v.(bool); ok {
		return vb
	}

	panic("not a bool")
}

func rString(v interface{}) string {
	if v == nil {
		return ""
	} else if vf, ok := v.(float64); ok {
		return fmt.Sprintf("%f", vf)
	} else if vi, ok := v.(int); ok {
		return strconv.Itoa(vi)
	} else if vs, ok := v.(string); ok {
		return vs
	} else if vb, ok := v.(bool); ok {
		if vb {
			return "Yes"
		}
		return "No"
	}

	panic("not a bool")
}
