doc:
	godoc -http=:6060 # http://localhost:6060/pkg/k/



transfer-to-saas:
	rsync -av --exclude "go.mod" --exclude "go.sum" --exclude ".cursor"  --exclude ".idea" --exclude "Makefile" ./  /Users/mtgnorton/Coding/go/src/github.com/mtgnorton/api-wikitrade-saas/common/k/
