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

	w := worker.New(client_var, model.PlaceOrderTaskQueue, worker.Options{})
	w.RegisterWorkflow(work.WorkFlow.RegisterWorkflow)
	w.RegisterWorkflow(work.WorkFlow.ViewDriverOrdersWorkflow)
	w.RegisterWorkflow(work.WorkFlow.OrderPlacedWorkflow)
	// w.RegisterActivity(work.Act.RegisterCheckRole)
	// w.RegisterActivity(work.Act.CreateUser)
	w.RegisterActivity(work.Act.ViewOrders)
	w.RegisterActivity(work.Act.GetItems)
	w.RegisterActivity(work.Act.CalculateBill)
	w.RegisterActivity(work.Act.SendEmail)
	w.RegisterActivity(work.Act.CreateOrder)

	log.Print("worker started")
	// Start listening to the Task Queues
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}

}
