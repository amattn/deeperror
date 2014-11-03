package deeperror

const (
	internal_BUILD_NUMBER   = 34
	internal_VERSION_STRING = "1.1.1"
)

func BuildNumber() int64 {
	return internal_BUILD_NUMBER
}
func Version() string {
	return internal_VERSION_STRING
}
