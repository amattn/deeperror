package deeperror

const (
	internalBuildNumber   = 3
	internalVersionString = "1.2.0"
)

func BuildNumber() int64 {
	return internalBuildNumber
}
func Version() string {
	return internalVersionString
}
