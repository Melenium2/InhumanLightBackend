package sqlstore

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/store"
)

type NotificationRepository struct {
	store *Store
}

func (repo *NotificationRepository) Create(newModel *models.Notification) error {
	if err := newModel.Validate(); err != nil {
		return err
	}
	newModel.BeforeCreate()

	return repo.store.db.QueryRow(
		"insert into notifications (mess, created_at, noti_status, for_user, checked) values($1, $2, $3, $4, $5) returning id",
		newModel.Message,
		newModel.Date,
		newModel.Status,
		newModel.For,
		newModel.Checked,
	).Scan(&newModel.ID)
}

func (repo *NotificationRepository)	FindById(userId uint) ([]*models.Notification, error) {
	rows, err := repo.store.db.Query(
		"select * from notifications where for_user = $1 and checked = false",
		userId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}
	defer rows.Close()

	notifications := make([]*models.Notification, 0)
	for rows.Next() {
		notification := &models.Notification{}
		if err := rows.Scan(&notification.ID, &notification.Message, &notification.Date,
		&notification.Status, &notification.For, &notification.Checked); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (repo *NotificationRepository)	Check(indexes []int, userId uint) error {
	sIndexes := make([]string, 0)
	for _, n := range indexes {
		str := strconv.Itoa(n)
		sIndexes = append(sIndexes, str)
	}
	
	_, err := repo.store.db.Exec(
		fmt.Sprintf("update notifications set checked = true where id in (%s) and for_user = %d", strings.Join(sIndexes, ","), userId),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}
		return err
	}

	return nil
}