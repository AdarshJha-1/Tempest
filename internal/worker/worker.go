package worker

import (
	"encoding/json"
	"log"

	"github.com/AdarshJha-1/Tempest/internal/executor"
	"github.com/AdarshJha-1/Tempest/internal/queue"
	"github.com/AdarshJha-1/Tempest/internal/store"
	"github.com/AdarshJha-1/Tempest/internal/types"
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
		var cfg types.Config
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

	}
}
