/* this is db record
*/

package base

type RecordWrapper interface {
	Parse(rcd *Record) error
	Build() (*Record, error)
}

type Record struct {
	names []string
	values []interface{}
}

func NewRecord() (*Record, error) {
	rcd := new(Record)
	err := rcd.init()
	return rcd, err
}

func (rcd *Record) init() error {
	rcd.names = make([]string, 0)
	rcd.values = make([]interface{}, 0)
	return nil
}

func (rcd *Record) Value(name string) interface{} {
	for i, v := range(rcd.names) {
		if v == name {
			return rcd.values[i]
		}
	}
	return nil
}

func (rcd *Record) SetValue(name string, value interface{}) {
	for i, v := range(rcd.names) {
		if v == name {
			rcd.values[i] = value
			return
		}
	}
	rcd.names = append(rcd.names, name)
	rcd.values = append(rcd.values, value)
}

func (rcd *Record) Names() []string {
	return rcd.names
}

func (rcd *Record) Values() []interface{} {
	return rcd.values
}