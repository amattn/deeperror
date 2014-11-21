package deeperror

const (
	internal_BUILD_NUMBER   = 35
	internal_VERSION_STRING = "1.1.2b"
)

func BuildNumber() int64 {
	return internal_BUILD_NUMBER
}
func Version() string {
	return internal_VERSION_STRING
}
