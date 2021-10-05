package model

const (
	userStatusDeactivate = -1
	userStatusRegistered = 0
	userStatusActive     = 1

	DefaultThumbnail = "https://s4.anilist.co/file/anilistcdn/character/large/b1336-73LQxWKUWy78.png"
)

type User struct {
	UID       string `json:"uid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Thumbnail string `json:"thumbnail"`
}

// func (u *User) getUserInfo() error {
// 	var (
// 		err   error
// 		query string
// 	)

// 	query = `SELECT first_name, last_name, thumbnail FROM user_info WHERE uid = $1 AND status = $2`
// 	if err = DB.QueryRowContext(context.Background(), query, u.UID, userStatusActive).Scan(&u.FirstName, u.LastName, u.Thumbnail); err != nil {
// 		log.Printf("[getUserInfo][QueryRowContext] Input: %s Output: %v", u.UID, err)
// 		return err
// 	}

// 	return err
// }

// func (u *User) addUserInfo(status int) error {
// 	var (
// 		err   error
// 		query string
// 	)

// 	u.Thumbnail = DefaultThumbnail

// 	query = `
// 		INSERT INTO user_info
// 		(uid, first_name, last_name, thumbnail, create_time, status)
// 		VALUES
// 		($1, $2, $3, $4, $5, $6)
// 	`
// 	if _, err = DB.ExecContext(context.Background(), query, u.UID, u.FirstName, u.LastName, u.Thumbnail, time.Now().Format(time.RFC3339), status); err != nil {
// 		log.Printf("[addUserInfo][ExecContext] Input: %s Output: %v", u.UID, err)
// 		return err
// 	}

// 	return err
// }
