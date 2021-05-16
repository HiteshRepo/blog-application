module github.com/HiteshRepo/blog-application

go 1.15

require (
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/improbable-eng/grpc-web v0.14.0
	github.com/rs/cors v1.7.0 // indirect
	github.com/stretchr/testify v1.7.0
	go.mongodb.org/mongo-driver v1.5.1
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b
	golang.org/x/net v0.0.0-20210510120150-4163338589ed
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
	nhooyr.io/websocket v1.8.7 // indirect
)

replace gopkg.in/urfave/cli.v2 => github.com/urfave/cli/v2 v2.3.0
