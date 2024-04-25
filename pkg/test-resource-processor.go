package pkg

import (
	"github.com/apibrew/apibrew/pkg/api"
	"github.com/apibrew/apibrew/pkg/errors"
	model2 "github.com/apibrew/template/pkg/model"
)

type testResourceProcessor struct {
	api                    api.Interface
	oauth2ConfigRepository api.Repository[*model2.TestResource]
}

func (t testResourceProcessor) Mapper() Mapper[*model2.TestResource] {
	return model2.TestResourceMapperInstance
}

func (t testResourceProcessor) Register(entity *model2.TestResource) error {
	entity.Description = entity.Name

	return nil
}

func (t testResourceProcessor) Update(entity *model2.TestResource) error {
	entity.Description = entity.Name

	return nil
}

func (t testResourceProcessor) UnRegister(entity *model2.TestResource) error {
	return errors.LogicalError.WithMessage("unregister not implemented")
}
