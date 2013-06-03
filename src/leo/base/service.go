/* this is service
*/

package base

type Service interface {
	Start() error
	Close() error
	Tick() error
	Save() error
}