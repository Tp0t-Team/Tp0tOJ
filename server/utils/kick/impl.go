package kick

import "sync"

var bannedUserId = sync.Map{}

func BanUser(user uint64) {
	bannedUserId.Store(user, true)
}

func UnbanUser(user uint64) {
	bannedUserId.Delete(user)
}

func KickGuard(user uint64) bool {
	_, ok := bannedUserId.Load(user)
	return !ok
}
