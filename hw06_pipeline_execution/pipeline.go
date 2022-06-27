package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func runStage(stage Stage, in In, done In) Out {
	myOut := make(Bi)

	go func() {
		defer close(myOut)

		stageOut := stage(in)

		for {
			select {
			case <-done:
				return
			case v, ok := <-stageOut:
				if !ok {
					return
				}
				myOut <- v
			}
		}
	}()

	return myOut
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	currentInput := in
	for _, stage := range stages {
		currentInput = runStage(stage, currentInput, done)
	}

	return currentInput
}
