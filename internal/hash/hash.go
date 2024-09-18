package hash

type Hash interface {
	Create(password string) (string, error)
	Compare(password, hash string) bool
}
