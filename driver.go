package dufu

type Driver interface {
	Probe()
	Remote()
}
