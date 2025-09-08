package waroom

import gen "github.com/topotal/waroom-client-go/internal/ogen"

// Generateした型から Exposeする型だけ書いていく
type Incident = gen.Incident
type IncidentLabel = gen.IncidentLabel
type IncidentMetrics = gen.IncidentMetrics
type IncidentSeverity = gen.IncidentSeverity

// TODO: 追加
