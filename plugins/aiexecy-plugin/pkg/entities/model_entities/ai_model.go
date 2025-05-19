package model_entities

import "aiexec-plugin/pkg/entities/plugin_entities"

type GetModelSchemasResponse struct {
	ModelSchema *plugin_entities.ModelDeclaration `json:"model_schema" validate:"omitempty"`
}
