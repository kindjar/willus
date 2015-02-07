package willus

type WillusError struct {
	Message string
}

func (e WillusError) Error() (error string) {
	return e.Message
}
