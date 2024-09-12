package cache

var defCreator = NewCacheCreator()

func GetCache(name string) ICache {
	return defCreator.GetCache(name)
}

func GetDefCache() ICache {
	return defCreator.GetDefCache()
}
