/* this is db record
*/

package base

type RecordWrapper interface {
	Parse(rcd *Record) error
	Build() (*Record, error)
}

type Record struct {
	Names []string
	Values []interface{}
}

func NewRecord() (*Record, error) {
	rcd := new(Record)
	err := rcd.init()
	return rcd, err
}

func (rcd *Record) init() error {
	rcd.Names = make([]string, 0)
	rcd.Values = make([]interface{}, 0)
	return nil
}

func (rcd *Record) Value(name string) interface{} {
	for i, v := range(rcd.Names) {
		if v == name {
			return rcd.Values[i]
		}
	}
	return nil
}

func (rcd *Record) SetValue(name string, value interface{}) {
	for i, v := range(rcd.Names) {
		if v == name {
			rcd.Values[i] = value
			return
		}
	}
	rcd.Names = append(rcd.Names, name)
	rcd.Values = append(rcd.Values, value)
}