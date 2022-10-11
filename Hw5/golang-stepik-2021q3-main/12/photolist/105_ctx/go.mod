module photolist

go 1.13

require (
	github.com/99designs/gqlgen v0.10.1
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/aws/aws-sdk-go v1.25.31
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/disintegration/imaging v1.6.1
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/golang/protobuf v1.3.2
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749
	github.com/spf13/viper v1.5.0
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/vektah/gqlparser v1.1.2
	golang.org/x/crypto v0.0.0-20191029031824-8986dd9e96cf
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sys v0.0.0-20190801041406-cbf593c0f2f3 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/grpc v1.23.0
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.25.1
