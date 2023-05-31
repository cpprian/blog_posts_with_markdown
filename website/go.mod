module github.com/cpprian/blog_posts_with_markdown/website

go 1.20

require (
	github.com/cpprian/blog_posts_with_markdown/users v0.0.0-20230528140956-b443575c61a7
	github.com/gomarkdown/markdown v0.0.0-20230322041520-c84983bdbf2a
	github.com/gorilla/mux v1.8.0
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
)

require github.com/cpprian/blog_posts_with_markdown/comments v0.0.0-20230531093234-9cdd4c61bf47

require (
	github.com/cpprian/blog_posts_with_markdown/posts v0.0.0-20230531095057-131ab7422b13
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	go.mongodb.org/mongo-driver v1.11.6 // indirect
)
