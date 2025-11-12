module github.com/rockcookies/go-caches/providers/redka

go 1.23.0

require (
	github.com/nalgeon/redka v0.6.0
	github.com/rockcookies/go-caches v0.0.0
)

require github.com/mattn/go-sqlite3 v1.14.32 // indirect

replace github.com/rockcookies/go-caches => ../..
