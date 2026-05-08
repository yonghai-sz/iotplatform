
`goctl model mysql ddl -c -src shorturl.sql -dir .`




`mockgen -destination ./shorturlmodel_mock.go -package model . ShorturlModel`  
`mockgen -destination shorturlmodel_mock.go -package model shorturlmodel_gen.go ShorturlModel`

