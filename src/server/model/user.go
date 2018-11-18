package model

import (
    "database/sql"
    "fmt"
)

/**
 * DESCRIPTION:
 *
 * @author rd
 * @create 2018-11-08 23:15
 **/

type User struct {
    ID      int64  `db:"id"`
    Account string `db:"account"`
    Phone   string `db:"phone"`
    Pass    string `db:"pass"`
}

func (*User) TableName() string {
    return "api_user"
}

func (u *User) GetAsAccount(db *sql.DB, account string) (*User, error) {
    queryStr := fmt.Sprintf("SELECT id, account, phone, pass FROM %s WHERE account=$1", u.TableName())
    rows, err := db.Query(queryStr, account)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var res *User
    for rows.Next() {
        var id sql.NullInt64
        var acc, ph, pass sql.NullString
        if err = rows.Scan(&id, &acc, &ph, &pass); err != nil {
            return nil, err
        }
        res = &User{
            ID:      id.Int64,
            Account: acc.String,
            Phone:   ph.String,
            Pass:    pass.String,
        }
    }
    return res, nil
}
