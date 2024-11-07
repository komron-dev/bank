# Ongoing pet project about bank (not completed)

A simple banking service that offers APIs to support:
	1.	Creating and managing bank accounts, each including the account holder’s name, balance, and currency type.
	2.	Logging all balance updates for every account, ensuring that each deposit or withdrawal generates a new record.
	3.	Facilitating money transfers between two accounts, using transactions to ensure that both accounts are updated successfully or, if an issue occurs, no changes are made.

## Setup local development

### Install tools

- [Docker Desktop](https://www.docker.com/products/docker-desktop)
- [PopSQL](https://popsql.com/) - free for users with Github Student Developer Pack
- [Golang](https://golang.org/)

### MacOS 
```bash
brew install golang-migrate
```

### Windows Using [scoop](https://scoop.sh/) 
```bash 
scoop install migrate
``` 
### Linux (*.deb package) 
```bash 
curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
apt-get update
apt-get install -y migrate
```


- [Sqlc](https://github.com/kyleconroy/sqlc#installation)

    ```bash
    brew install sqlc
    ```

- [Gomock](https://github.com/golang/mock)

    ``` bash
    go install github.com/golang/mock/mockgen
    ```

### Setup infrastructure

- Create the bank-network for Docker

    ``` bash
    make network
    ```

- Start postgres container in Docker:

    ```bash
    make postgres
    ```

- Create simplebank database:

    ```bash
    make create_db
    ```

- Run db migration:

    ```bash
    make up_migrate
    ```

- Run db migration up 1 version:

    ```bash
    make up_migrate_last
    ```

- Run db migration down all versions:

    ```bash
    make down_migrate
    ```

- Run db migration down 1 version:

    ```bash
    make down_migrate_last
    ```

- Generate SQL CRUD with sqlc:

    ```bash
    make sqlc
    ```

- Generate DB mock with gomock:

    ```bash
    make mock
    ```

- Create a new db migration:

    ```bash
    make new_migration name=<migration_name>
    ```

### How to run

- Run the program:

    ```bash
    make run
    ```

- Run test:

    ```bash
    make test
    ```
