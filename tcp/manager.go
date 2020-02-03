package tcp

type manager struct{}

func newSessionManager() *manager {
	return &manager{}
}
