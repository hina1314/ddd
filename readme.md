## DEPLOYMENT
* Installing sqlc [doc](https://docs.sqlc.dev/en/latest/overview/install.html#installing-sqlc), requires Go 1.21+.
  > go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
* Installing golang-migrate [github](https://github.com/golang-migrate/migrate) , change $database to your database  
  > go install -tags '$database' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
* Run
  > make migrate_init
* then write your db schema in generated files in ./db/migration/.
* Run 
  > make migrate_up

## Customize Your Model
* you can write your queries in ./db/query/ folder.
* for example, refer to [sqlc](https://docs.sqlc.dev/en/latest/howto/select.html)
  ```sql 
  -- name: GetUserByID :one
  SELECT * FROM "user" WHERE id = $1;
  ```
* Then run
  > make sqlc