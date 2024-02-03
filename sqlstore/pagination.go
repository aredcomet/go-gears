package sqlstore

import (
	"errors"
	"github.com/aredcomet/go-gears/utils"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
)

var decoder = schema.NewDecoder()

type PageParams struct {
	Limit  int    `schema:"limit"`
	Order  string `schema:"order"`
	Cursor string `schema:"cursor"`
}

func (p *PageParams) AsPaginator(paginatorFunc func(paginator.Cursor, *paginator.Order, *int) *paginator.Paginator) *paginator.Paginator {
	cursor := paginator.Cursor{After: &p.Cursor}
	limit := p.Limit
	order := paginator.Order(p.Order)
	if paginatorFunc == nil {
		return CreateDefaultPaginator(cursor, &order, &limit)
	}
	return paginatorFunc(cursor, &order, &limit)
}

func GetPageParamOrReject(w http.ResponseWriter, r *http.Request, logger *logrus.Logger) (PageParams, error) {
	var params PageParams
	err := decoder.Decode(&params, r.URL.Query())
	if err != nil {
		msg := "Unable to parse URL parameters"
		utils.RespondWithError(w, http.StatusBadRequest, msg, logger)
		return params, errors.New(msg)

	}
	if params.Order != "ASC" && params.Order != "DESC" && params.Order != "" {
		msg := "order should be ASC or DESC"
		utils.RespondWithError(w, http.StatusBadRequest, msg, logger)
		return params, errors.New(msg)
	}

	if params.Limit == 0 || params.Limit > 100 {
		params.Limit = 20
	}
	return params, nil
}

func CreateDefaultPaginator(
	cursor paginator.Cursor,
	order *paginator.Order,
	limit *int,
) *paginator.Paginator {
	opts := []paginator.Option{
		&paginator.Config{
			Keys:  []string{"ID"},
			Limit: 20,
			Order: paginator.ASC,
		},
	}
	if limit != nil {
		opts = append(opts, paginator.WithLimit(*limit))
	}
	if order != nil {
		opts = append(opts, paginator.WithOrder(*order))
	}
	if cursor.After != nil {
		opts = append(opts, paginator.WithAfter(*cursor.After))
	}
	if cursor.Before != nil {
		opts = append(opts, paginator.WithBefore(*cursor.Before))
	}
	return paginator.New(opts...)
}
