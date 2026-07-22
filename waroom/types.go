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

// 新specはcomponents/schemasを持たずレスポンススキーマが全てインライン定義のため、
// 一覧と単体取得で別の生成型になる
type Incident = gen.GetIncidentOK
type IncidentListItem = gen.ListIncidentsOKIncidentsItem
type IncidentLabel = gen.GetIncidentOKLabelsItem
type IncidentMetrics = gen.GetIncidentOKMetrics
type ListIncidentsParams = gen.ListIncidentsParams

type Label = gen.ListServiceLabelsOKLabelsItem
type CreatedLabel = gen.CreateServiceLabelCreatedLabel
type UpdatedLabel = gen.UpdateServiceLabelOKLabel
type Service = gen.ListServicesOKServicesItem

// ListServicesParams は新specでクエリパラメータが無くなったため空
type ListServicesParams struct{}
type ListServiceLabelsParams = gen.ListServiceLabelsParams

type CreateServiceLabelParams = gen.CreateServiceLabelParams
type CreateServiceLabelReq = gen.CreateServiceLabelReqLabel

type UpdateServiceLabelParams = gen.UpdateServiceLabelParams
type UpdateServiceLabelReq = gen.UpdateServiceLabelReqLabel
type DeleteServiceLabelParams = gen.DeleteServiceLabelParams

type User = gen.ListUsersOKUsersItem
type ListUsersParams = gen.ListUsersParams

type Runbook = gen.GetRunbookOKRunbook
type RunbookListItem = gen.ListRunbooksOKRunbooksItem
type CreatedRunbook = gen.CreateRunbookCreatedRunbook
type UpdatedRunbook = gen.UpdateRunbookOKRunbook
type CreateRunbookReq = gen.CreateRunbookReqRunbook
type UpdateRunbookReq = gen.UpdateRunbookReqRunbook
type ListRunbooksParams = gen.ListRunbooksParams
type GetRunbookParams = gen.GetRunbookParams
type UpdateRunbookParams = gen.UpdateRunbookParams
type DeleteRunbookParams = gen.DeleteRunbookParams

type RunbookRuleListItem = gen.ListServiceRunbookRulesOKRunbookRulesItem
type CreatedRunbookRule = gen.CreateServiceRunbookRuleCreated
type UpdatedRunbookRule = gen.UpdateServiceRunbookRuleOK
type CreateRunbookRuleReq = gen.CreateServiceRunbookRuleReqRunbookRule
type UpdateRunbookRuleReq = gen.UpdateServiceRunbookRuleReqRunbookRule
type ListServiceRunbookRulesParams = gen.ListServiceRunbookRulesParams
type CreateServiceRunbookRuleParams = gen.CreateServiceRunbookRuleParams
type UpdateServiceRunbookRuleParams = gen.UpdateServiceRunbookRuleParams
type DeleteServiceRunbookRuleParams = gen.DeleteServiceRunbookRuleParams

type WebhookConfigurationListItem = gen.ListServiceWebhookConfigurationsOKWebhookConfigurationsItem
type WebhookConfiguration = gen.GetServiceWebhookConfigurationOKWebhookConfiguration
type ListServiceWebhookConfigurationsParams = gen.ListServiceWebhookConfigurationsParams
type GetServiceWebhookConfigurationParams = gen.GetServiceWebhookConfigurationParams

type PostmortemListItem = gen.ListPostmortemsOKPostmortemsItem
type Postmortem = gen.GetPostmortemOKPostmortem
type CreatedPostmortem = gen.CreatePostmortemOK
type UpdatedPostmortem = gen.UpdatePostmortemOKPostmortem
type CreatePostmortemReq = gen.CreatePostmortemReq
type UpdatePostmortemReq = gen.UpdatePostmortemReq
type PostmortemDefaultTemplate = gen.GetPostmortemTemplateOK
type UpdatedPostmortemDefaultTemplate = gen.UpdatePostmortemTemplateOK
type UpdatePostmortemTemplateReq = gen.UpdatePostmortemTemplateReq
type ListPostmortemsParams = gen.ListPostmortemsParams
type GetPostmortemParams = gen.GetPostmortemParams
type UpdatePostmortemParams = gen.UpdatePostmortemParams
type DeletePostmortemParams = gen.DeletePostmortemParams

type PostmortemTemplateListItem = gen.ListPostmortemTemplatesOKPostmortemTemplatesItem
type PostmortemTemplate = gen.GetPostmortemTemplateByUuidOKPostmortemTemplate
type CreatedPostmortemTemplate = gen.CreatePostmortemTemplateCreatedPostmortemTemplate
type UpdatedPostmortemTemplate = gen.UpdatePostmortemTemplateByUuidOKPostmortemTemplate
type CreatePostmortemTemplateReq = gen.CreatePostmortemTemplateReqPostmortemTemplate
type UpdatePostmortemTemplateByUuidReq = gen.UpdatePostmortemTemplateByUuidReqPostmortemTemplate
type ListPostmortemTemplatesParams = gen.ListPostmortemTemplatesParams
type GetPostmortemTemplateByUuidParams = gen.GetPostmortemTemplateByUuidParams
type UpdatePostmortemTemplateByUuidParams = gen.UpdatePostmortemTemplateByUuidParams
type DeletePostmortemTemplateByUuidParams = gen.DeletePostmortemTemplateByUuidParams

type ActionItemListItem = gen.ListIncidentActionItemsOKActionItemsItem
type ActionItem = gen.GetIncidentActionItemOKActionItem
type CreatedActionItem = gen.CreateIncidentActionItemCreatedActionItem
type CreateActionItemReq = gen.CreateIncidentActionItemReqActionItem
type ListIncidentActionItemsParams = gen.ListIncidentActionItemsParams
type CreateIncidentActionItemParams = gen.CreateIncidentActionItemParams
type GetIncidentActionItemParams = gen.GetIncidentActionItemParams
type DeleteIncidentActionItemParams = gen.DeleteIncidentActionItemParams

type IncidentDescriptionTemplate = gen.GetIncidentDescriptionTemplateOK
type UpdatedIncidentDescriptionTemplate = gen.UpdateIncidentDescriptionTemplateOK
type UpdateIncidentDescriptionTemplateReq = gen.UpdateIncidentDescriptionTemplateReqIncidentDescriptionTemplate

// TODO: 追加
