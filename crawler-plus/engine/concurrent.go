package engine

//ConcurrentEngine
//开启并发爬虫采集器
type ConcurrentEngine struct {
	Scheduler   Scheduler
	Fetcher     Fetcher
	WorkerCount int
	Writer      WriteWorker
}

//Run
func (c *ConcurrentEngine) Run(seeds ...Request) {
	var out = make(chan ParseResult)
	c.Scheduler.Begin()
	for i := 0; i < c.WorkerCount; i++ {
		CreateWorkers(out, c.Scheduler, c.Fetcher)
	}
	for _, request := range seeds {
		c.Scheduler.Submit(request)
	}
	for {
		result := <-out
		for _, item := range result.Items {
			go func() { c.Writer.Payload <- item }()
		}
		for _, request := range result.Requests {
			c.Scheduler.Submit(request)
		}

	}
}

//CreateWorkers
func CreateWorkers(out chan ParseResult, s Scheduler, f Fetcher) {
	in := make(chan Request)
	go func() {
		for {
			s.WorkChanFree(in)
			requests := <-in
			result, err := f.Work(requests)
			if err != nil {
				continue
			}
			out <- result
		}

	}()
}
