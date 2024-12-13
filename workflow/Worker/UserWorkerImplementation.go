package worker

import (
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func (work *Worker) WorkerUserStart() {
	client_var, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer client_var.Close()

	w := worker.New(client_var, model.RegisterTaskQueue, worker.Options{})
	w.RegisterWorkflow(work.WorkFlow.RegisterWorkflow)
	w.RegisterActivity(work.Act.RegisterCheckRole)
	w.RegisterActivity(work.Act.CreateUser)
	log.Print("worker started")
	// Start listening to the Task Queuess
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}

}
