package scheduler

import (
	"context"
	"time"

	"github.com/mankings/mec-federator/internal/models/dto"
)

// CreateTasks initializes all periodic tasks
func CreateTasks(scheduler *Scheduler) {
	scheduler.AddTask(Task{
		Name:     "FederationAppis",
		Interval: 10 * time.Second,
		Handler: func(ctx context.Context) error {
			// get app instances from federator db
			appInstances, err := scheduler.services.AppInstanceService.GetAllAppInstances()
			if err != nil {
				return err
			}

			// convert into list of appi ids
			appiIds := make([]string, 0)
			for _, appInstance := range appInstances {
				appiIds = append(appiIds, appInstance.Id)
			}

			// get app instances from orchestrator db
			orchAppis, err := scheduler.services.OrchestratorService.GetAppInstances(appiIds)
			if err != nil {
				return err
			}

			// make message to send to kafka
			var message dto.FederatedAppisMessage
			message.Appis = make(map[string]dto.AppiDetails)
			for _, appi := range orchAppis {
				message.Appis[appi.AppiId] = dto.AppiDetails{
					Domain:    appi.Domain,
					Instances: appi.Instances[appi.Domain],
				}
			}

			// send message to kafka
			_, err = scheduler.services.KafkaClientService.Produce("federation-appis", message)
			if err != nil {
				return err
			}

			return err
		},
	})

	// Example: Cleanup old messages from Kafka service
	scheduler.AddTask(Task{
		Name:     "KafkaMessageCleanup",
		Interval: 10 * time.Minute,
		Handler: func(ctx context.Context) error {
			scheduler.services.KafkaClientService.CleanupOldMessages(10 * time.Minute)
			return nil
		},
	})
}
