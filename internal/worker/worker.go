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

		resp, err := w.executor.Run(&cfg)
		if err != nil {
			continue
		}

		if resp {
			w.store.UpdateJobStatusById(jobId, "success")
			w.store.UpdateJobFinishTimeById(jobId)
		}
	}
}
