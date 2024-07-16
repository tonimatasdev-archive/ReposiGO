package session

var sessions = make(map[string]Session)

type Storage struct {
}
