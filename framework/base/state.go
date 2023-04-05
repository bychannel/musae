package base

type PState int32

const (
	PState_None     PState = 0
	PState_Starting PState = 1 // starting
	PState_Running  PState = 2 // running
	PState_Loading  PState = 3 // loading
	PState_Exiting  PState = 4 // exiting
)
