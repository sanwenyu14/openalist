package meilisearch

import (
	"errors"
	"fmt"
	"github.com/AlliotTech/openalist/internal/conf"
	"github.com/AlliotTech/openalist/internal/model"
	"github.com/AlliotTech/openalist/internal/search/searcher"
	"github.com/AlliotTech/openalist/pkg/utils"
	"github.com/meilisearch/meilisearch-go"
)

var config = searcher.Config{
	Name:       "meilisearch",
	AutoUpdate: true,
}

func init() {
	searcher.RegisterSearcher(config, func() (searcher.Searcher, error) {
		m := Meilisearch{
			Client: meilisearch.NewClient(meilisearch.ClientConfig{
				Host:   conf.Conf.Meilisearch.Host,
				APIKey: conf.Conf.Meilisearch.APIKey,
			}),
			IndexUid:             conf.Conf.Meilisearch.IndexPrefix + "alist",
			FilterableAttributes: []string{"parent", "is_dir", "name"},
			SearchableAttributes: []string{"name"},
		}

		_, err := m.Client.GetIndex(m.IndexUid)
		if err != nil {
			var mErr *meilisearch.Error
			ok := errors.As(err, &mErr)
			if ok && mErr.MeilisearchApiError.Code == "index_not_found" {
				task, err := m.Client.CreateIndex(&meilisearch.IndexConfig{
					Uid:        m.IndexUid,
					PrimaryKey: "id",
				})
				if err != nil {
					return nil, err
				}
				forTask, err := m.Client.WaitForTask(task.TaskUID)
				if err != nil {
					return nil, err
				}
				if forTask.Status != meilisearch.TaskStatusSucceeded {
					return nil, fmt.Errorf("index creation failed, task status is %s", forTask.Status)
				}
			} else {
				return nil, err
			}
		}
		attributes, err := m.Client.Index(m.IndexUid).GetFilterableAttributes()
		if err != nil {
			return nil, err
		}
		if attributes == nil || !utils.SliceAllContains(*attributes, m.FilterableAttributes...) {
			_, err = m.Client.Index(m.IndexUid).UpdateFilterableAttributes(&m.FilterableAttributes)
			if err != nil {
				return nil, err
			}
		}

		attributes, err = m.Client.Index(m.IndexUid).GetSearchableAttributes()
		if err != nil {
			return nil, err
		}
		if attributes == nil || !utils.SliceAllContains(*attributes, m.SearchableAttributes...) {
			_, err = m.Client.Index(m.IndexUid).UpdateSearchableAttributes(&m.SearchableAttributes)
			if err != nil {
				return nil, err
			}
		}

		pagination, err := m.Client.Index(m.IndexUid).GetPagination()
		if err != nil {
			return nil, err
		}
		if pagination.MaxTotalHits != int64(model.MaxInt) {
			_, err := m.Client.Index(m.IndexUid).UpdatePagination(&meilisearch.Pagination{
				MaxTotalHits: int64(model.MaxInt),
			})
			if err != nil {
				return nil, err
			}
		}
		return &m, nil
	})
}
