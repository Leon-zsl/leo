/* this is component base
*/

package comp

type Component interface {
	ID() string
	Type() string
	Start() error
	Close() error
	Tick() error
}