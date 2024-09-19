package ports

var CacheContext *cacheContext = &cacheContext{}

type cacheContext struct {
	CheckEmailCache       Cache
	CodeConfirmationCache Cache
}

type Cache interface {
	Set(key string, v interface{}) error
	Get(key string, v interface{}) error
	GetAndDelete(key string, v interface{}) error
	GetString(key string) (string, error)
	GetStringAndDelete(key string) (string, error)
}
