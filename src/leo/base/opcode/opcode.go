/* this is the opcode
*/

package opcode

const (
	//sys opcode
	LOGIN = 1 << 16
	
	//player opcode
	GETINFO = 2 << 16

	//battle opcode
	CAST_SKILL = 3 << 16

)