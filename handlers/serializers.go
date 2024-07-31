package handlers

import "flowspell/models"

func serializeTaskDefinition(td models.TaskDefinition) TaskDefinitionResponse {
    return TaskDefinitionResponse{
        ID:                  td.ID,
        CreatedAt:           td.CreatedAt,
        UpdatedAt:           td.UpdatedAt,
        ReferenceID:         td.ReferenceID,
        Name:                td.Name,
        Description:         td.Description,
        FlowDefinitionRefID: td.FlowDefinitionRefID,
        InputSchema:         td.InputSchema,
        OutputSchema:        td.OutputSchema,
        Version:             td.Version,
        Metadata:            td.Metadata,
    }
}

func serializeFlowDefinition(fd models.FlowDefinition) FlowDefinitionResponse {
    return FlowDefinitionResponse{
        ID:           fd.ID,
        ReferenceID:  fd.ReferenceID,
        CreatedAt:    fd.CreatedAt,
        UpdatedAt:    fd.UpdatedAt,
        Name:         fd.Name,
        Description:  fd.Description,
        Status:       fd.Status,
        Version:      fd.Version,
        InputSchema:  fd.InputSchema,
        OutputSchema: fd.OutputSchema,
        Metadata:     fd.Metadata,
    }
}
