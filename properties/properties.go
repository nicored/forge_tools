package properties

import (
	"encoding/json"
	"path/filepath"
)

const (
	tBool    = 1
	tColor   = 2
	tNumeric = 3
	tObjRef  = 11
	tString  = 20
	tString2 = 21
)

const (
	aDefParent           = "__parent__/parent"
	aDefName             = "__name__/name"
	aDefInstOf           = "__instanceof__/instanceof_objid"
	aDefViewableIn       = "__viewable_in__/viewable_in"
	aDefChild            = "__child__/child"
	aDefNodeFlags        = "__node_flags__/node_flags"
	aDefDocSchemaName    = "__document__/schema_name"
	aDefDocSchemaVersion = "__document__/schema_version"
	aDefIsDocProperty    = "__document__/is_doc_property"
)

type Properties struct {
	attrs Attrs
	offs  Offs
	avs   Avs
	vals  Vals
	ids   Ids

	dirPath string
}

func NewProperties(path string) *Properties {
	p := &Properties{
		attrs: readAttrs(filepath.Join(path, "objects_attrs.json")),
		offs:  readOffs(filepath.Join(path, "objects_offs.json")),
		avs:   readAvs(filepath.Join(path, "objects_avs.json")),
		vals:  readVals(filepath.Join(path, "objects_vals.json")),
		ids:   readIds(filepath.Join(path, "objects_ids.json")),

		dirPath: path,
	}

	return p
}

type ObjectProps struct {
	ObjectID   int
	Name       string
	ExternalID string
	Properties map[string]map[string]interface{}
	Parents    []int
}

func (p *Properties) ExportJson(indent bool) ([]byte, error) {
	results := p.Run()

	if indent {
		return json.MarshalIndent(results, "    ", "  ")
	}

	return json.Marshal(results)
}

func (p *Properties) Run() []*ObjectProps {
	var dbIds []int
	for i := 1; i <= p.offs.IdMax(); i++ {
		dbIds = append(dbIds, i)
	}

	var objs []*ObjectProps
	for _, dbId := range dbIds {
		if obj := p.getObjectProperties(dbId); obj != nil {
			objs = append(objs, obj)
		}
	}

	return objs
}

func (p *Properties) getObjectProperties(dbId int) *ObjectProps {
	result := &ObjectProps{
		ObjectID:   dbId,
		Properties: map[string]map[string]interface{}{},
		ExternalID: p.ids[dbId],
		Parents:    []int{},
	}

	parent := p.read(dbId, result)

	for parent >= 0 && parent != 1 {
		parent = p.read(parent, result)
	}

	if result.Name == "" {
		return nil
	}

	return result
}

func (p *Properties) read(dbId int, result *ObjectProps) int {
	parent := -1
	propStart := 2 * p.offs[dbId]
	propStop := len(p.avs)

	if len(p.offs) > dbId+1 {
		propStop = 2 * p.offs[dbId+1]
	}

	for i := propStart; i < propStop; i += 2 {
		attr := p.attrs[p.avs[i]-1]

		key := attr.Category + "/" + attr.Name

		if key == aDefParent {
			parent = int(p.vals[p.avs[i+1]].(float64))
			result.Parents = append(result.Parents, parent)
			continue
		}

		if key == aDefInstOf {
			p.read(int(p.vals[p.avs[i+1]].(float64)), result)
			continue
		}

		if key == aDefViewableIn ||
			key == aDefChild ||
			key == aDefNodeFlags ||
			key == aDefDocSchemaName ||
			key == aDefDocSchemaVersion ||
			key == aDefIsDocProperty {
			continue
		}

		if key == aDefName && (*attr.Type == tString || *attr.Type == tString2) {
			if result.Name == "" {
				result.Name = p.vals[p.avs[i+1]].(string)
			}

			continue
		}

		if _, ok := result.Properties[attr.Category]; !ok {
			result.Properties[attr.Category] = map[string]interface{}{}
		}

		result.Properties[attr.Category][attr.Name] = parsePropValue(attr, p.vals[p.avs[i+1]])
	}

	return parent
}

func parsePropValue(attr *Attr, v interface{}) string {
	value := rString(v)

	if attr.Unit != nil {
		value += " " + *attr.Unit
	}

	return value
}
