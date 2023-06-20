package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = runner(done, stage(in))
	}
	return in
}

func runner(done In, in In) Bi {
	out := make(Bi)

	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case value, ok := <-in:
				if !ok {
					return
				}
				out <- value
			}
		}
	}()

	return out
}
