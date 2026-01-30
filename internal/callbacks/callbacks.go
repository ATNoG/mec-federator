package callbacks

import (
	"context"

	"github.com/atnog/mec-federator/internal/router"
)

func StartCallbacks(services *router.Services) {
	kafkaServ := services.KafkaClientService
	kafkaServ.StartConsumer(context.Background(), "responses", nil, false) // use default callback for responses, might change in the future

	infrastructureInfoCallback := NewInfrastructureInfoCallback(services)
	kafkaServ.StartConsumer(context.Background(), "infrastructure-info", infrastructureInfoCallback.HandleMessage, false)

	federationInfrastructureInfoCallback := NewFederationInfrastructureInfoCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation-infrastructure-info", federationInfrastructureInfoCallback.HandleMessage, false)

	getAppiInfoCallback := NewGetAppiInfoCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation_get_appi_info", getAppiInfoCallback.HandleMessage, true)

	newFederationCallback := NewNewFederationCallback(services)
	kafkaServ.StartConsumer(context.Background(), "new_federation", newFederationCallback.HandleMessage, true)

	removeFederationCallback := NewRemoveFederationCallback(services)
	kafkaServ.StartConsumer(context.Background(), "remove_federation", removeFederationCallback.HandleMessage, true)

	getFederationInfoCallback := NewGetFederationInfoCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation_get_info", getFederationInfoCallback.HandleMessage, true)

	federationZoneSubscribeCallback := NewFederationZoneSubscribeCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation_subscribe_zone", federationZoneSubscribeCallback.HandleMessage, true)

	federationZoneUnsubscribeCallback := NewFederationZoneUnsubscribeCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation_unsubscribe_zone", federationZoneUnsubscribeCallback.HandleMessage, true)

	newFederationArtefactCallback := NewFederationArtefactNewCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation_new_artefact", newFederationArtefactCallback.HandleMessage, true)

	removeFederationArtefactCallback := NewFederationArtefactRemoveCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation_remove_artefact", removeFederationArtefactCallback.HandleMessage, true)

	newFederationAppCallback := NewFederationNewAppCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation_new_app", newFederationAppCallback.HandleMessage, true)

	removeFederationAppCallback := NewFederationRemoveAppCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation_remove_app", removeFederationAppCallback.HandleMessage, true)

	newFederationAppiCallback := NewFederationAppiNewCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation_new_appi", newFederationAppiCallback.HandleMessage, true)

	removeFederationAppiCallback := NewFederationRemoveAppiCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation_remove_appi", removeFederationAppiCallback.HandleMessage, true)

	enableAppInstKDUCallback := NewEnableAppInstanceKDUCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation_enable_kdu", enableAppInstKDUCallback.HandleMessage, true)

	disableAppInstKDUCallback := NewDisableAppInstanceKDUCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation_disable_kdu", disableAppInstKDUCallback.HandleMessage, true)

	migrateAppInstNodeCallback := NewFederationMigrateNodeCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation_migrate_node", migrateAppInstNodeCallback.HandleMessage, true)

}
