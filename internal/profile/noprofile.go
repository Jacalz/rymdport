//go:build !profile

package profile

func Start() func() {
	return nil
}
