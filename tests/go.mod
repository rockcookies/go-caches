module github.com/rockcookies/go-caches/tests

go 1.23.0

require (
	github.com/mattn/go-sqlite3 v1.14.32
	github.com/nalgeon/redka v0.6.0
	github.com/redis/go-redis/v9 v9.16.0
	github.com/rockcookies/go-caches v0.0.0
	github.com/rockcookies/go-caches/providers/redis v0.0.0
	github.com/rockcookies/go-caches/providers/redka v0.0.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/rockcookies/go-caches => ../
	github.com/rockcookies/go-caches/providers/redis => ../providers/redis
	github.com/rockcookies/go-caches/providers/redka => ../providers/redka
)
