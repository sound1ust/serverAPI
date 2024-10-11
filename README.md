### This is my Go API
:)

### Dependencies
- [Migrate CLI](https://github.com/golang-migrate/migrate/tree/v4.18.1/cmd/migrate)

### Running the project
Run your PostgreSQL, then create a database and run the migrations:

```bash
make migrate-up
```

After that, you can run the project with the following command:

```bash
make run
```