package domain

type User struct {
	Name   string
	PwHash string
}

type Users []User

func (ul Users) CheckUser(name string) (exists bool, pwhash string) {
	for _, val := range ul {
		if val.Name == name {
			return true, val.PwHash
		}
	}
	return false, ""
}
