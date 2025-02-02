package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in, done In, stages ...Stage) Out {
	if in == nil {
		res := make(Bi)
		close(res)
		return res
	}
	out := in
	for _, stage := range stages {
		out = stage(processStage(out, done))
	}

	return out
}

func processStage(in, done In) Out {
	processOut := make(Bi)
	go func() {
		defer func() {
			close(processOut)
			for range in {
				continue
			}
		}()
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case processOut <- v:
				}
			}
		}
	}()
	return processOut
}
