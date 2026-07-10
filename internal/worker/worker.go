package worker

import (
	"encoding/json"
	"log"

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

func (w *worker) Start() {

	// currently my each job is processed one by one
	for {
		jobId, err := w.queue.GetJobID()
		if err != nil {
			log.Println(err)
			continue
		}

		configBytes, err := w.store.GetJobConfigById(jobId)
		if err != nil {
			continue
		}
		var cfg config.Config
		err = json.Unmarshal(configBytes, &cfg)
		if err != nil {
			continue
		}

		w.store.UpdateJobStatusById(jobId, "running")

		result, err := w.executor.Run(&cfg)
		if err != nil {
			w.store.UpdateJobStatusById(jobId, "failed")
			w.store.UpdateJobFinishTimeById(jobId)
			continue
		}

		// here i need to save result
		err = w.store.InsertResult(jobId, result)
		if err != nil {
			// idk what to do here
			log.Println("failed to insert result in db", err)
		}

		err = w.store.UpdateJobStatusById(jobId, "completed")
		if err != nil {
			// idk what to do here
			log.Println("failed to update job status in db", err)
		}

		err = w.store.UpdateJobFinishTimeById(jobId)
		if err != nil {
			// idk what to do here
			log.Println("failed to update job finish time in db", err)
		}
	}
}
