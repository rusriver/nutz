package executionlog

type ExecutionLog struct {
	Records []*Record
}

type Record struct {
	InstructionId string
	ResultString  string
}
