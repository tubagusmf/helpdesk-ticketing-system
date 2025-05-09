package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"helpdesk-ticketing-system/internal/model"

	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
)

type TicketHistoryRepo struct {
	db       *gorm.DB
	esClient *elastic.Client
}

func NewTicketHistoryRepo(db *gorm.DB, esClient *elastic.Client) model.ITicketHistoryRepository {
	return &TicketHistoryRepo{
		db:       db,
		esClient: esClient,
	}
}

func (t *TicketHistoryRepo) IndexToElasticsearch(history *model.TicketHistory) error {
	ctx := context.Background()

	_, err := t.esClient.Index().
		Index("ticket_history").
		Id(fmt.Sprintf("%d", history.ID)).
		BodyJson(history).
		Do(ctx)

	return err
}

func (t *TicketHistoryRepo) GetTicketID(ctx context.Context, id int64) (*model.TicketHistory, error) {
	var history model.TicketHistory

	err := t.db.WithContext(ctx).
		Where("ticket_id = ?", id).
		Order("changed_at DESC").
		First(&history).Error
	if err != nil {
		return nil, err
	}

	return &history, err
}

func (t *TicketHistoryRepo) GetStatus(ctx context.Context, status string) (*[]model.TicketHistory, error) {
	var histories []model.TicketHistory

	query := elastic.NewMatchQuery("status", status)

	searchResult, err := t.esClient.Search().
		Index("ticket_history").
		Query(query).
		Sort("changed_at", false).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	for _, hit := range searchResult.Hits.Hits {
		var history model.TicketHistory
		err := json.Unmarshal(hit.Source, &history)
		if err != nil {
			return nil, err
		}

		histories = append(histories, history)
	}

	return &histories, err
}

func (t *TicketHistoryRepo) GetPriority(ctx context.Context, priority string) (*[]model.TicketHistory, error) {
	var histories []model.TicketHistory

	query := elastic.NewMatchQuery("priority", priority)

	searchResult, err := t.esClient.Search().
		Index("ticket_history").
		Query(query).
		Sort("changed_at", false).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	for _, hit := range searchResult.Hits.Hits {
		var history model.TicketHistory
		err := json.Unmarshal(hit.Source, &history)
		if err != nil {
			return nil, err
		}

		histories = append(histories, history)
	}

	return &histories, err
}

func (t *TicketHistoryRepo) GetUserID(ctx context.Context, userID int64) (*[]model.TicketHistory, error) {
	var histories []model.TicketHistory

	query := elastic.NewMatchQuery("user_id", userID)

	searchResult, err := t.esClient.Search().
		Index("ticket_history").
		Query(query).
		Sort("changed_at", false).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	for _, hit := range searchResult.Hits.Hits {
		var history model.TicketHistory
		err := json.Unmarshal(hit.Source, &history)
		if err != nil {
			return nil, err
		}

		histories = append(histories, history)
	}

	return &histories, err
}

func (t *TicketHistoryRepo) Create(ctx context.Context, ticketHistory model.TicketHistory) error {
	err := t.db.WithContext(ctx).Create(&ticketHistory).Error
	if err != nil {
		return err
	}

	_, err = t.esClient.Index().
		Index("ticket_history").
		Id(fmt.Sprintf("%d", ticketHistory.ID)).
		BodyJson(ticketHistory).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to index ticket history to Elasticsearch: %w", err)
	}

	return nil
}
