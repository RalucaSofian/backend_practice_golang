# Project Name: backend_practice_golang
## Project Description

The aim of the project is to create the back-end for a project representing a Pet Rescue,
where Users can register, search for a Pet, and adopt or foster it.


## Technologies Used

- GoLang
- Chi Router
- Bun
- PostgreSQL


## Project Folder Structure

Source files of the project are saved under the following folder structure:

```bash
.
└── app/
    ├── controllers/
    ├── db/
    ├── infra/
    │   └── docker-compose.yml
    ├── migrations/
    ├── models/
    ├── services/
    ├── types/
    ├── utils/
    │   ├── access_utils/
    │   ├── date/
    │   ├── middlewares/
    │   ├── query_utils/
    │   └── ...
    ├── cli.go
    └── main.go
```


## Installing the Project

The following commands are used for installing and running the project:
```bash
# start the local database
$: cd ./app/infra
$: docker compose up -d
✔ Container infra-local_db_go-1  Started
```
```bash
# create your .env file:
# LISTEN_PORT=****
# DB_HOSTNAME=****
# DB_PORT=****
# DB_USERNAME=****
# DB_PASSWORD=****
# DB_NAME=****
# SECRET_KEY=****

# start the local (development) server
$: go run ./app
```

Upon successful start of the development server, the following log will be printed:
```bash
[server] API Server Running on Port ...
```


## Functionalities

- AUTH:
    - Register: password hashing and salting
    - Login: Bearer Token
- CRUD:
    - Individual GET by ID
    - Querying:
        - Filtering
        - Searching
        - Ordering
        - Pagination
    - Individual UPDATE and DELETE of entities
- Permission system based on User Roles
- GoLang migration system


## Migration System

The migration system makes use of the GoLang Bun database client. The migration commands can be listed by running the following:
```bash
$: go run ./app help
```

The list of available commands will be printed:
```bash
[CLI] Available CLI commands

[CLI] create-migration : Create Migration with given name ; Parameters: migration name: string
[CLI] apply-migrations : Apply all unapplied Migrations ; Parameters: none
[CLI] revert-migrations : Revert most recently applied Migrations ; Parameters: none
[CLI] list-migrations : List all Migrations ; Parameters: none
[CLI] help : CLI Migrations Help ; Parameters: none
```

As an example, in order to create and apply (or revert) a migration, the following commands will be used:
```bash
$: go run ./app create-migration name_of_migration_file
[CLI] Created Migration: 20241210130120_name_of_migration_file.go

$: go run ./app list-migrations
[CLI] List of Unapplied Migrations is: 20241210130120_name_of_migration_file

$: go run ./app apply-migrations
 [up migration] [CLI] Applied Migrations: group #6 (20241210130120_name_of_migration_file)

$: go run ./app revert-migrations
 [down migration] [CLI] Reverted Migrations: group #6 (20241210130120_name_of_migration_file)
```
