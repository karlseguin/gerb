package core

type ExecutionState int
type TagType int

const (
	NormalState ExecutionState = iota

	OutputTag TagType = iota
	UnsafeTag
	CodeTag
	NoTag
)
