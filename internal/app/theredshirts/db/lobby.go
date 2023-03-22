package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

const (
	lobby_table_name       = "lobby"
	create_lobby_sql       = "INSERT INTO %s.%s(id, name, owner, password, difficulty) VALUES($1, $2, $3, $4, $5)"
	update_lobby_sql       = "UPDATE %s.%s SET name = $2, owner = $3, password = $4, difficulty = $5) WHERE id = $1"
	delete_lobby_sql       = "DELETE FROM %s.%s WHERE id = $1"
	select_lobby_by_id_sql = "SELECT id, name, owner, password, difficulty FROM %s.%s WHERE id = $1"
	select_lobby_sql       = "SELECT id, name, owner, password, difficulty FROM %s.%s"
)

var (
	ErrLobbyAlreadyExists = errors.New("lobby already exists")
)

func (db *postgresConnection) CreateLobby(lobby *Lobby) error {
	if _, err := db.dbPool.Exec(context.Background(), fmt.Sprintf(create_lobby_sql, schema_name, lobby_table_name), lobby.ID, lobby.Name, lobby.Password, lobby.Difficulty); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return ErrLobbyAlreadyExists
			}
		}

		return fmt.Errorf("unknown error when inserting lobby: %v", err)
	}
	return nil
}

func (db *postgresConnection) UpdateLobby(lobby *Lobby) error {
	if _, err := db.dbPool.Exec(context.Background(), fmt.Sprintf(update_lobby_sql, schema_name, lobby_table_name), lobby.ID, lobby.Name, lobby.Password, lobby.Difficulty); err != nil {
		return fmt.Errorf("unknown error when updating lobby: %v", err)
	}
	return nil
}

func (db *postgresConnection) DeleteLobby(id uuid.UUID) error {
	if _, err := db.dbPool.Exec(context.Background(), fmt.Sprintf(delete_lobby_sql, schema_name, lobby_table_name), id); err != nil {
		return fmt.Errorf("unknown error when deliting lobby: %v", err)
	}
	return nil
}

func (db *postgresConnection) GetLobbyById(id uuid.UUID) (*Lobby, error) {
	var lobbies []*Lobby
	if err := pgxscan.Select(context.Background(), db.dbPool, &lobbies, fmt.Sprintf(select_lobby_by_id_sql, schema_name, lobby_table_name), id); err != nil {
		return nil, fmt.Errorf("error while selecting lobby with id %v: %v", id, err)
	}

	if len(lobbies) != 1 {
		return nil, fmt.Errorf("cant find only one lobby. Lobbies: %v", lobbies)
	}

	return lobbies[0], nil
}

func (db *postgresConnection) GetAllLobbies() ([]*Lobby, error) {
	var lobbies []*Lobby
	if err := pgxscan.Select(context.Background(), db.dbPool, &lobbies, fmt.Sprintf(select_lobby_sql, schema_name, lobby_table_name)); err != nil {
		return nil, fmt.Errorf("error while selecting all lobbies: %v", err)
	}

	return lobbies, nil
}