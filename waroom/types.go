package waroom

import (
	"time"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

// Generateした型から Exposeする型だけ書いていく

// General
type OptString = gen.OptString
type OptDateTime = gen.OptDateTime
type OptInt = gen.OptInt
type OptNilInt = gen.OptNilInt

func NewOptString(v string) OptString        { return gen.NewOptString(v) }
func NewOptDateTime(v time.Time) OptDateTime { return gen.NewOptDateTime(v) }
func NewOptInt(v int) OptInt                 { return gen.NewOptInt(v) }
func NewOptNillInt(v int) OptNilInt          { return gen.NewOptNilInt(v) }

type Incident = gen.Incident
type IncidentLabel = gen.IncidentLabel
type IncidentMetrics = gen.IncidentMetrics
type IncidentSeverity = gen.IncidentSeverity
type ListIncidentsParams = gen.ListIncidentsParams

type Label = gen.Label
type Service = gen.Service
type ListServicesParams = gen.ListServicesParams
type ListServiceLabelsParams = gen.ListServiceLabelsParams
type CreateServiceLabelParams = gen.CreateServiceLabelParams
type UpdateServiceLabelParams = gen.UpdateServiceLabelParams
type DeleteServiceLabelParams = gen.DeleteServiceLabelParams

// TODO: 追加
