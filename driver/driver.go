package driver

type Driver interface{
	// IdTables()
	// Flags()
	Probe()
	Remove()
}
