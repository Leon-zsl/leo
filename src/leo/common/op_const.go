/* this is opcode consts
*/

package common

const (
	//sys opcode
	OP_LOGIN = 1 << 16
	
	//player opcode
	OP_GETINFO = 2 << 16

	//battle opcode
	OP_CAST_SKILL = 3 << 16
)