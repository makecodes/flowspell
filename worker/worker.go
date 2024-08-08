package worker

import (
	fsTasks "flowspell/tasks"
	"os"

	"github.com/RichardKnop/machinery/v2"
	amqpbackend "github.com/RichardKnop/machinery/v2/backends/amqp"
	amqpbroker "github.com/RichardKnop/machinery/v2/brokers/amqp"
	"github.com/RichardKnop/machinery/v2/config"
	"github.com/RichardKnop/machinery/v2/example/tracers"
	eagerlock "github.com/RichardKnop/machinery/v2/locks/eager"
	"github.com/RichardKnop/machinery/v2/log"
	"github.com/RichardKnop/machinery/v2/tasks"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func startServer() (*machinery.Server, error) {
	brokerUrl := os.Getenv("BROKER_URL")
	defaultQueue := os.Getenv("DEFAULT_QUEUE")
	defaultExchange := os.Getenv("DEFAULT_EXCHANGE")
	defaultBindingKey := os.Getenv("DEFAULT_BINDING_KEY")
	cnf := &config.Config{
		Broker:          brokerUrl,
		DefaultQueue:    defaultQueue,
		ResultBackend:   brokerUrl,
		ResultsExpireIn: 3600,
		AMQP: &config.AMQPConfig{
			Exchange:      defaultExchange,
			ExchangeType:  "direct",
			BindingKey:    defaultBindingKey,
			PrefetchCount: 3,
		},
	}

	// Create server instance
	broker := amqpbroker.New(cnf)
	backend := amqpbackend.New(cnf)
	lock := eagerlock.New()
	server := machinery.NewServer(cnf, broker, backend, lock)

	// Register tasks
	tasksMap := map[string]interface{}{
		"queue": fsTasks.Queue,
	}

	return server, server.RegisterTasks(tasksMap)
}

func Worker() error {
	consumerTag := "flowspell_consumer_tag"

	cleanup, err := tracers.SetupTracer(consumerTag)
	if err != nil {
		log.FATAL.Fatalln("Unable to instantiate a tracer:", err)
	}
	defer cleanup()

	server, err := startServer()
	if err != nil {
		return err
	}

	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	worker := server.NewWorker(consumerTag, 0)

	// Here we inject some custom code for error handling,
	// start and end of task hooks, useful for metrics for example.
	errorHandler := func(err error) {
		log.ERROR.Println("I am an error handler:", err)
	}

	preTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("I am a start of task handler for:", signature.Name)
	}

	postTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("I am an end of task handler for:", signature.Name)
	}

	worker.SetPostTaskHandler(postTaskHandler)
	worker.SetErrorHandler(errorHandler)
	worker.SetPreTaskHandler(preTaskHandler)

	return worker.Launch()
}
