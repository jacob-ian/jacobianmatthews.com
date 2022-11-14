# jacobianmatthews.com/api

The backend for my website. It uses HTTP for Auth and GRPC for the CMS.

## Development

### Database State

To create a valid development database, SSH into the `backend` container and run:

- `make migration_up`
- `make seed_dev`

### Database Migrations

To run/create database migrations in development, you must shell into the `backend` development docker container.

- `make migration_create NAME=[NAME]` - Creates an `up` and `down` migration with the inputted name
- `make migration_up` - Runs the `up` migrations locally
- `make migration_down` - Runs the `down` migrations locally
