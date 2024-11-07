package oreillyapi

import (
	"fmt"
	"net/url"
)

const (
	QueryKey                    = "query"
	FieldKey                    = "field"
	FormatsKey                  = "formats"
	LanguagesKey                = "languages"
	Publisherskey               = "publishers"
	TopicUuidsKey               = "topic_uuids"
	TopicsKey                   = "topics"
	VideoClassificationsKey     = "video_classifications"
	CertificationVendorNamesKey = "certification_vendor_name"
	CertificationNamesKey       = "certification_name"
	SortKey                     = "sort"
	OrderKey                    = "order"
	PageKey                     = "page"
	LimitKey                    = "limit"
)

type SearchOption struct {
	Query                    string
	Field                    string
	Formats                  []string
	Languages                []string
	Publishers               []string
	TopicUuids               []string
	Topics                   []string
	VideoClassifications     []string
	CertificationVendorNames []string
	CertificationNames       []string
	Sort                     Sort
	Order                    Order
	Page                     int
	Limit                    int
}

type Sort string

const (
	Relevance       Sort = "relevance"
	Popularity      Sort = "popularity"
	Rating          Sort = "rating"
	DateAdded       Sort = "date_added"
	PublicationDate Sort = "publication_date"
)

type Order string

const (
	Asc  Order = "asc"
	Desc Order = "desc"
)

type Filters struct {
	Publishers []string
}

func (so *SearchOption) queryParams() url.Values {
	q := make(url.Values)

	// set query
	if so.Query == "" {
		q.Set(QueryKey, "*")
	} else {
		q.Set(QueryKey, so.Query)
	}

	// set field
	if so.Field != "" {
		q.Set(FieldKey, so.Field)
	}

	// set filters
	for _, v := range so.Formats {
		q.Add(FormatsKey, v)
	}
	for _, v := range so.Languages {
		q.Add(LanguagesKey, v)
	}
	for _, v := range so.Publishers {
		q.Add(Publisherskey, v)
	}
	for _, v := range so.TopicUuids {
		q.Add(TopicUuidsKey, v)
	}
	for _, v := range so.Topics {
		q.Add(TopicsKey, v)
	}
	for _, v := range so.VideoClassifications {
		q.Add(VideoClassificationsKey, v)
	}
	for _, v := range so.CertificationVendorNames {
		q.Add(CertificationVendorNamesKey, v)
	}
	for _, v := range so.CertificationNames {
		q.Add(CertificationNamesKey, v)
	}

	// set sorting
	if so.Sort != "" {
		q.Set(SortKey, string(so.Sort))
	}
	if so.Order != "" {
		q.Set(OrderKey, string(so.Order))
	}

	// set page, limit
	if so.Page > 0 {
		q.Set(PageKey, fmt.Sprintf("%d", so.Page))
	}
	if so.Limit > 0 {
		q.Set(LimitKey, fmt.Sprintf("%d", so.Limit))
	}

	return q
}
