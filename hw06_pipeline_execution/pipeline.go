package hw06pipelineexecution

import "runtime"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

const concurrentStages = 4 // consider using system parameter

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	outCh := startProducers(in, done, stages)

	out := startConsumer(done, outCh)

	return out
}

func startProducers(in In, done In, stages []Stage) chan Out {
	outCh := make(chan Out, runtime.NumCPU()*concurrentStages)
	go func() {
		defer close(outCh)
		for {
			select {
			case <-done:
				return
			case data, ok := <-in:
				if ok {
					select {
					case <-done:
						return
					case outCh <- performStages(startOnceProducer(done, data), stages...):
						// push into result wait channel
					}
				} else {
					return
				}
			}
		}
	}()
	return outCh
}

func startOnceProducer(done In, data interface{}) Out {
	dataIn := make(Bi)
	go func() {
		defer close(dataIn)
		select {
		case <-done:
		case dataIn <- data:
		}
	}()
	return dataIn
}

func performStages(in In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = stage(out)
	}
	return out
}

func startConsumer(done In, outCh chan Out) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case resCh, ok := <-outCh:
				if ok {
					receiveResult(resCh, done, out)
				} else {
					return
				}
			}
		}
	}()
	return out
}

func receiveResult(res Out, done In, out Bi) {
	for {
		select {
		case <-done:
			return
		case result, ok := <-res:
			if ok {
				select {
				case <-done:
					return
				case out <- result:
					// write data
				}
			} else {
				return
			}
		}
	}
}
