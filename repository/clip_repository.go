package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ritikmayank/catch-it/model"
)

type ClipRepository struct {
	db *pgxpool.Pool
}

func NewClipRepository(db *pgxpool.Pool) *ClipRepository {
	return &ClipRepository{db: db}
}

func (r *ClipRepository) ListByDate(ctx context.Context, date time.Time) ([]model.Clip, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, clip_index, text, language FROM clips
		 WHERE play_date = $1 ORDER BY clip_index`,
		date.Format("2006-01-02"),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clips []model.Clip
	for rows.Next() {
		var c model.Clip
		if err := rows.Scan(&c.ID, &c.ClipIndex, &c.Text, &c.Language); err != nil {
			return nil, err
		}
		clips = append(clips, c)
	}
	return clips, rows.Err()
}

func (r *ClipRepository) GetByID(ctx context.Context, id int) (*model.Clip, error) {
	var c model.Clip
	err := r.db.QueryRow(ctx,
		`SELECT id, play_date, clip_index, text, language FROM clips WHERE id=$1`, id,
	).Scan(&c.ID, &c.PlayDate, &c.ClipIndex, &c.Text, &c.Language)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ClipRepository) Upsert(ctx context.Context, date string, index int, text, language string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO clips (play_date, clip_index, text, language)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (play_date, clip_index) DO NOTHING`,
		date, index, text, language,
	)
	return err
}
