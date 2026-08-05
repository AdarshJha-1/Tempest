package worker

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/AdarshJha-1/Tempest/internal/config"
	"github.com/AdarshJha-1/Tempest/internal/executor"
	"github.com/AdarshJha-1/Tempest/internal/queue"
	"github.com/AdarshJha-1/Tempest/internal/store"
)

type Worker interface {
	Start()
}

type worker struct {
	queue     queue.Queue
	store     store.Store
	executor  executor.Executor
	workerCap int
}

func New(que queue.Queue, store store.Store, executor executor.Executor, workerCap int) Worker {
	return &worker{
		queue:     que,
		store:     store,
		executor:  executor,
		workerCap: workerCap,
	}
}

func (w *worker) processJob(jobId string) error {
	configBytes, err := w.store.GetJobConfig(jobId)
	if err != nil {
		return err
	}
	var cfg config.Config
	err = json.Unmarshal(configBytes, &cfg)
	if err != nil {
		return err
	}

	w.store.UpdateJobStatus(jobId, "running")

	result, err := w.executor.Run(&cfg)
	if err != nil {
		w.store.UpdateJobStatus(jobId, "failed")
		w.store.UpdateJobFinishTime(jobId)
		return err
	}

	err = w.store.CreateResult(jobId, result)
	if err != nil {
		log.Println("failed to insert result in db", err)
		return err
	}

	err = w.store.UpdateJobStatus(jobId, "completed")
	if err != nil {
		log.Println("failed to update job status in db", err)
		return err
	}

	err = w.store.UpdateJobFinishTime(jobId)
	if err != nil {
		log.Println("failed to update job finish time in db", err)
		return err
	}
	return nil
}

func (w *worker) run(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		jobId, err := w.queue.GetJobID()
		if err != nil {
			log.Println("error getting job from queue")
			continue
		}
		if err := w.processJob(jobId); err != nil {
			log.Println(err)
			continue
		}
	}
}

func (w *worker) Start() {
	wg := sync.WaitGroup{}

	for i := 0; i < w.workerCap; i++ {
		wg.Add(1)
		go w.run(&wg)
	}
	wg.Wait()
}
