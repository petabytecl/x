package jobx

type Job interface {
	Process()
}

type FuncExecutorJob struct {
	Err  error
	Func func() error
}

func (j *FuncExecutorJob) Process() {
	j.Err = j.Func()
}

var _ Job = &FuncExecutorJob{}
