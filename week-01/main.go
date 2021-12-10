package main

import (
    "context"
    "database/sql"
    "log"

    "github.com/pkg/errors"
)

var (
    ctx context.Context
    db  *sql.DB
)

// Select 根据用户id查询用户昵称，使用 errors.Wrapf 包装错误
func Select(id int) (string ,error) {
    var username string
    err := db.QueryRowContext(ctx, "SELECT username FROM users WHERE id=?", id).Scan(&username)
    return username, errors.Wrapf(err, "SELECT username FROM users WHERE id=%d", id)
}

func main() {
    uid := 123
    username, err := Select(uid) // 查询用户id
    if errors.Is(err, sql.ErrNoRows) { // 空行
        log.Printf("no such user: %v", err)
        return
    }
    if err != nil { // 其他的错误
        log.Fatalf("query error: %v\n", err)
        return
    }
    // 正常逻辑
    log.Printf("id: %d, username: %s\n", uid, username)
}
