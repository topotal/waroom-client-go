package waroom

import gen "github.com/topotal/waroom-client-go/internal/ogen"

// Generateした型から Exposeする型だけ書いていく
type Incident = gen.Incident
type IncidentLabel = gen.IncidentLabel
type IncidentMetrics = gen.IncidentMetrics
type IncidentSeverity = gen.IncidentSeverity
type ListIncidentsParams = gen.ListIncidentsParams

// TODO: services

type Label = gen.Label
type ListServicesParams = gen.ListServicesParams

// TODO: 追加
