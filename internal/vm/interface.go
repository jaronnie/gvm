package vm

type Interface interface {
	List() ([]string, error)
}
