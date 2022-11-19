package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = startStage(out, done, stage)
	}
	return out
}

func startStage(in In, done In, stage Stage) Out {
	intermediate := make(Bi)

	go func() {
		defer close(intermediate)
		for {
			select {
			case <-done:
				return
			case data, ok := <-in:
				if ok {
					select {
					case <-done:
						return
					case intermediate <- data:
						// write to channel
					}
				} else {
					return
				}
			}
		}
	}()

	return stage(intermediate)
}
