package repositories

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type PortiereUserInfo struct {
	ID      string `db:"id"`
	AliasID int64  `db:"alias_id"`
}

func GetLegacyIDs(limit int) (ids []string, err error) {
	rdb := getHermesInstance()
	statement := `SELECT user_id FROM user_push_token WHERE alias_id IS NULL LIMIT ?`

	query, args, err := sqlx.In(statement, limit)
	if err != nil {
		log.Println("GetLegacyIDs err:", err)
		return nil, err
	}
	query = rdb.Rebind(query)
	return ids, rdb.Select(&ids, query, args...)
}

func GetPortiereUserInfo(legacyIDs []string) (result []PortiereUserInfo, err error) {
	rdb := getPortiereMasterInstance()
	statement := `SELECT
					users.id AS id,
					aliases.id AS alias_id
				FROM
					aliases
					JOIN users ON users.snowflake_user_id = aliases.user_id
				WHERE
					aliases.alias_type = 1
					AND users.id IN (?)`

	query, args, err := sqlx.In(statement, legacyIDs)
	if err != nil {
		log.Println("GetPortiereUserInfo err:", err)
		return nil, err
	}
	query = rdb.Rebind(query)
	return result, rdb.Select(&result, query, args...)
}

func MigrationPortiereUsersInfoToHermes(users []PortiereUserInfo) {
	rdb := getHermesInstance()

	statement := `UPDATE user_push_token SET alias_id = $1 WHERE user_id = $2`
	query := rdb.Rebind(statement)

	for _, element := range users {
		if _, err := rdb.Exec(query, element.AliasID, element.ID); err != nil {
			log.Println("MigrationPortiereUsersInfoToHermes Exec err:", err)
		}
	}
}