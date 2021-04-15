package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/woremacx/quartz"
)

//demo main
func main() {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go demoJobs(wg)
	go demoScheduler(wg)

	wg.Wait()
}

func demoScheduler(wg *sync.WaitGroup) {
	sched := quartz.NewStdScheduler()
	cronTrigger, _ := quartz.NewCronTrigger("1/3 * * * * *")
	cronJob := PrintJob{"Cron job"}
	sched.Start()
	sched.ScheduleJob(&PrintJob{"Ad hoc Job"}, quartz.NewRunOnceTrigger(time.Second*5))
	sched.ScheduleJob(&PrintJob{"First job"}, quartz.NewSimpleTrigger(time.Second*12))
	sched.ScheduleJob(&PrintJob{"Second job"}, quartz.NewSimpleTrigger(time.Second*6))
	sched.ScheduleJob(&PrintJob{"Third job"}, quartz.NewSimpleTrigger(time.Second*3))
	sched.ScheduleJob(&cronJob, cronTrigger)

	time.Sleep(time.Second * 10)

	j, _ := sched.GetScheduledJob(cronJob.Key())
	fmt.Println(j.TriggerDescription)
	fmt.Println("Before delete: ", sched.GetJobKeys())
	sched.DeleteJob(cronJob.Key())
	fmt.Println("After delete: ", sched.GetJobKeys())

	time.Sleep(time.Second * 2)
	sched.Stop()
	wg.Done()
}

func demoJobs(wg *sync.WaitGroup) {
	sched := quartz.NewStdScheduler()
	sched.Start()
	cronTrigger, _ := quartz.NewCronTrigger("1/5 * * * * *")
	shellJob := quartz.NewShellJob("ls -la")
	curlJob, err := quartz.NewCurlJob(http.MethodGet, "http://worldclockapi.com/api/json/est/now", "", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sched.ScheduleJob(shellJob, cronTrigger)
	sched.ScheduleJob(curlJob, quartz.NewSimpleTrigger(time.Second*7))

	time.Sleep(time.Second * 10)

	fmt.Println(sched.GetJobKeys())
	fmt.Println(shellJob.Result)
	fmt.Println(curlJob.Response)

	time.Sleep(time.Second * 2)
	sched.Stop()
	wg.Done()
}
