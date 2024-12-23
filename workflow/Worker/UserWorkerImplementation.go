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
	w.RegisterWorkflow(work.WorkFlow.PlaceOrderWorkflow)
	w.RegisterActivity(work.Act.FetchUserEmail)
	w.RegisterActivity(work.Act.GetItems)
	w.RegisterActivity(work.Act.CalculateBill)
	w.RegisterActivity(work.Act.SendEmail)
	w.RegisterActivity(work.Act.CreateOrder)
	w.RegisterActivity(work.Act.CheckOrderStatus)
	w.RegisterActivity(work.Act.CreateSourceConnection)
	w.RegisterActivity(work.Act.CreateDestinationConnection)
	w.RegisterActivity(work.Act.FetchSourceConfiguration)
	w.RegisterActivity(work.Act.FetchDestinationConfiguration)
	w.RegisterActivity(work.Act.MoveDataFromSourceToDestination)

	log.Print("worker started")
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}

}
