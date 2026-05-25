module github.com/clowz-me/catalog

go 1.25.5

replace github.com/clowz-me/pkg => ../pkg

require (
	github.com/clowz-me/pkg v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.3.0
	github.com/go-chi/cors v1.2.2
	github.com/google/uuid v1.6.0
)

require github.com/lib/pq v1.12.3 // indirect
