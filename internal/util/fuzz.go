package util

func FuzzCodeValidator(data []byte) int {
	validate := CodeValidator
	_ = validate(string(data))
	return 1
}
