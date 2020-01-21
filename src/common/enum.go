package common


// Define the GRPC server state in ENUM
type FizzbuzzServerState int
const (
	Init FizzbuzzServerState = 0
	Boot FizzbuzzServerState = 1
	Listen FizzbuzzServerState = 2
	Ready FizzbuzzServerState = 3
	Error FizzbuzzServerState = 4
	Gracefull FizzbuzzServerState = 5
	Stop FizzbuzzServerState = 6
)
