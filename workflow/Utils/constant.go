package utils

const Cancelled string = "cancelled"
const Accepted string = "accepted"
const Completed string = "completed"
const OrderPlaced string = "order placed"
const Delay string = "delay"

const CancelledSubject string = "Subject: Order Cancellation\n"
const AcceptedSubject string = "Subject: Order Confirmation\n"
const CompletedSubject string = "Subject: Order Completed\n"
const OrderPlacedSubject string = "Subject: Order placed\n"
const DelaySubject string = "Subject: Order Delayed\n"

const RegisterTaskQueue = "register_task_queue"
const PlaceOrderTaskQueue = "place_order_task_queue"
const DataSyncTaskQueue = "data_sync_task_queue"
