/* this is room module
*/

package module

type Module interface {
	ID() string
	Type() string
	Start() error
	Close() error
	Tick() error
}