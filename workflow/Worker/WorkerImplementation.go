package worker

import (
	"log"

	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func (work *Worker) WorkerUserStart() {
	client_var, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer client_var.Close()

	w := worker.New(client_var, utils.PlaceOrderTaskQueue, worker.Options{})
	w.RegisterWorkflow(work.WorkFlow.PlaceOrderWorkflow)
	w.RegisterActivity(work.Act.FetchUserEmail)
	w.RegisterActivity(work.Act.GetItems)
	w.RegisterActivity(work.Act.CalculateBill)
	w.RegisterActivity(work.Act.SendEmail)
	w.RegisterActivity(work.Act.CreateOrder)
	w.RegisterActivity(work.Act.CheckOrderStatus)

	w2 := worker.New(client_var, utils.DataSyncTaskQueue, worker.Options{})
	w2.RegisterWorkflow(work.WorkFlow.DataSyncWorkflow)
	w2.RegisterActivity(work.Act.CreateSourceToken)
	w2.RegisterActivity(work.Act.CreateDestinationToken)
	w2.RegisterActivity(work.Act.FetchSourceConfiguration)
	w2.RegisterActivity(work.Act.FetchDestinationConfiguration)
	w2.RegisterActivity(work.Act.MoveDataFromSourceToDestination)
	w2.RegisterActivity(work.Act.AddLogs)
	w2.RegisterActivity(work.Act.CountFilesInFolder)
	w2.RegisterActivity(work.Act.MoveBatchActivity)

	// log.Print("worker started")
	// err = w.Run(worker.InterruptCh())
	// if err != nil {
	// 	log.Fatalln("unable to start Worker", err)
	// }

	// Start both workers in separate goroutines
	go func() {
		log.Print("PlaceOrder worker started")
		err = w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalln("unable to start PlaceOrder worker", err)
		}
	}()

	go func() {
		log.Print("DataSync worker started")
		err = w2.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalln("unable to start DataSync worker", err)
		}
	}()

	// Block main thread until interruption
	select {}

}
