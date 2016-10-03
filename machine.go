package golorp

type Machine struct {
	// M0
	Code       []Cell
	Heap       []Cell
	PDL        []*Cell
	XRegisters []Cell
	HReg       int
	Mode       InstructionMode

	// M1
	ARegisters []Cell
	PReg       int

	// M2
	AndStack []Environment
	CPReg    int
	EReg     int

	// M3 - Prolog
	OrStack []ChoicePoint
	Trail   []Cell
	BReg    int
	TRReg   int
	GBReg   int

	// Optimisations
	B0Reg int
}

type Environment Cell
type ChoicePoint Cell

type Instruction int
type InstructionMode int

const (
	InvalidInstruction Instruction = iota // zero value not sure

	// L0
	PutStructure
	SetVariable
	SetValue
	GetStructure
	UnifyVariable
	UnifyValue

	// L1
	Call
	Proceeed
	PutVariable
	PutValue
	GetVariable
	GetValue

	// L2
	allocate
	deallocate

	// L3 - Prolog
	TryMeElse
	RetryMeElse
	TrustMe

	// Optimisations
)

const (
	InvalidMode InstructionMode = iota // zero value not sure
	Read
	Write
)
