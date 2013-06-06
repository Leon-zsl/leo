/* this is ai base
*/

package world

type AI interface {
	Start() error 
	Close() error
	Tick() error
}