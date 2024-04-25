package pkg

import (
	"github.com/apibrew/apibrew/pkg/api"
	"github.com/apibrew/apibrew/pkg/errors"
	"github.com/apibrew/apibrew/pkg/model"
	"github.com/apibrew/apibrew/pkg/resources"
	"github.com/apibrew/apibrew/pkg/service"
	"github.com/apibrew/apibrew/pkg/service/backend-event-handler"
	"github.com/apibrew/apibrew/pkg/util"
	model2 "github.com/apibrew/template/pkg/model"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
)

type module struct {
	container           service.Container
	backendEventHandler backend_event_handler.BackendEventHandler
	api                 api.Interface
}

func (m module) Init() {
	m.ensureNamespace()
	m.ensureResources()

	oauth2ConfigRepository := api.NewRepository[*model2.TestResource](m.api, model2.TestResourceMapperInstance)

	if err := RegisterResourceProcessor[*model2.TestResource](
		"sso-oauth2-authenticate-listener",
		&testResourceProcessor{
			api:                    m.api,
			oauth2ConfigRepository: oauth2ConfigRepository,
		},
		m.backendEventHandler,
		m.container,
		model2.TestResourceResource,
	); err != nil {
		log.Fatal(err)
	}

}

func (m module) ensureNamespace() {
	_, err := m.container.GetRecordService().Apply(util.SystemContext, service.RecordUpdateParams{
		Namespace: resources.NamespaceResource.Namespace,
		Resource:  resources.NamespaceResource.Name,
		Records: []*model.Record{
			{
				Properties: map[string]*structpb.Value{
					"name": structpb.NewStringValue("template"),
				},
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}

func (m module) ensureResources() {
	var list = []*model.Resource{
		model2.TestResourceResource,
	}

	for _, resource := range list {
		existingResource, err := m.container.GetResourceService().GetResourceByName(util.SystemContext, resource.Namespace, resource.Name)

		if err == nil {
			resource.Id = existingResource.Id
			err = m.container.GetResourceService().Update(util.SystemContext, resource, true, true)

			if err != nil {
				log.Fatal(err)
			}
		} else if err.Is(errors.ResourceNotFoundError) {
			_, err = m.container.GetResourceService().Create(util.SystemContext, resource, true, true)

			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}
}

func NewModule(container service.Container) service.Module {
	a := api.NewInterface(container)

	backendEventHandler := container.GetBackendEventHandler().(backend_event_handler.BackendEventHandler)
	return &module{container: container,
		api:                 a,
		backendEventHandler: backendEventHandler}
}
