package consts

const WorkerIDKey = "WorkerIDKey"

type ContextKey struct {
	key string
}

func NewContextKey(key string) *ContextKey {
	return &ContextKey{
		key: key,
	}
}

var WorkerIDContextKey = NewContextKey(WorkerIDKey)
