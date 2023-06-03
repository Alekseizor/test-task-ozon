package links

import (
	"context"
	"database/sql"
)

type RepoLinkPostgres struct {
	Ctx context.Context
	DB  *sql.DB
}

func NewRepoLinkPostgres(db *sql.DB, ctx context.Context) (*RepoLinkPostgres, error) {
	return &RepoLinkPostgres{
		DB:  db,
		Ctx: ctx,
	}, nil

}

func (lm *RepoLinkPostgres) AddLink(item *Links) error {
	_, err := lm.DB.ExecContext(lm.Ctx, "INSERT INTO link VALUES ($1,$2);", item.InitialURL, item.ShortenURL)
	if err != nil {
		return err
	}
	return nil
}

func (lm *RepoLinkPostgres) GetInitialLink(url string) (*Links, error) {
	row := lm.DB.QueryRowContext(lm.Ctx, "SELECT initial_url,shorten_url FROM link WHERE shorten_url=$1 LIMIT 1;", url)
	link := new(Links)
	err := row.Scan(&link.InitialURL, &link.ShortenURL)
	if err != nil {
		return nil, err
	}
	return link, nil
}

func (lm *RepoLinkPostgres) GetShortenLink(url string) (*Links, error) {
	row := lm.DB.QueryRowContext(lm.Ctx, "SELECT initial_url,shorten_url FROM link WHERE initial_url=$1 LIMIT 1;", url)
	link := new(Links)
	err := row.Scan(&link.InitialURL, &link.ShortenURL)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return link, nil
}
