package worker

type EmailDispatcher interface {
	Dispatch(job EmailJob)
}
