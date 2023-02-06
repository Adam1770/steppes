package client

type Pipeline interface {
	CheckPipeline() bool
	Source() error
	Build() error
	Test() error
	Deploy() error
}
