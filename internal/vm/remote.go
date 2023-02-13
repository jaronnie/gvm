package vm

type RemoteVM struct {
	Registry string
}

func NewRemoteVM() Interface {
	return &RemoteVM{}
}

func (o RemoteVM) List() ([]string, error) {
	return nil, nil
}
