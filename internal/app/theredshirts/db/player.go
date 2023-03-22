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
	player_table_name              = "player"
	create_player_sql              = "INSERT INTO %s.%s(id, name, lobby_id) VALUES($1, $2, $3)"
	delete_player_sql              = "DELETE FROM %s.%s WHERE id = $1"
	select_player_by_player_id_sql = "SELECT id, name, lobby_id FROM %s.%s WHERE id = $1"
	select_player_by_lobby_id_sql  = "SELECT id, name, lobby_id FROM %s.%s WHERE lobby_id = $1"
)

var (
	ErrPlayerAlreadyExists = errors.New("player already exists")
)

func (db *postgresConnection) CreatePlayer(player *Player) error {
	if _, err := db.dbPool.Exec(context.Background(), fmt.Sprintf(create_player_sql, schema_name, player_table_name), player.ID, player.Name, player.LobbyId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return ErrPlayerAlreadyExists
			}
		}

		return fmt.Errorf("unknown error when inserting player: %v", err)
	}
	return nil
}

func (db *postgresConnection) DeletePlayer(id uuid.UUID) error {
	if _, err := db.dbPool.Exec(context.Background(), fmt.Sprintf(delete_player_sql, schema_name, player_table_name), id); err != nil {
		return fmt.Errorf("unknown error when deliting player: %v", err)
	}
	return nil
}

func (db *postgresConnection) GetPlayerById(id uuid.UUID) (*Player, error) {
	var players []*Player
	if err := pgxscan.Select(context.Background(), db.dbPool, &players, fmt.Sprintf(select_player_by_player_id_sql, schema_name, player_table_name), id); err != nil {
		return nil, fmt.Errorf("error while selecting player with id %v: %v", id, err)
	}

	if len(players) != 1 {
		return nil, fmt.Errorf("cant find only one player. Players: %v", players)
	}

	return players[0], nil
}

func (db *postgresConnection) GetAllPlayersInLobby(lobbyId uuid.UUID) ([]*Player, error) {
	var players []*Player
	if err := pgxscan.Select(context.Background(), db.dbPool, &players, fmt.Sprintf(select_player_by_lobby_id_sql, schema_name, player_table_name), lobbyId); err != nil {
		return nil, fmt.Errorf("error while selecting all players in lobby: %v", err)
	}

	return players, nil
}