package deeperror

const (
	internal_BUILD_NUMBER   = 36
	internal_VERSION_STRING = "1.1.2"
)

func BuildNumber() int64 {
	return internal_BUILD_NUMBER
}
func Version() string {
	return internal_VERSION_STRING
}
