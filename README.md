# TyreMatch

TyreMatch is a computer vision backend system developed as part of an **MSc Forensic Investigation** dissertation research project. The system explores the feasibility of using image-based computer vision techniques to identify vehicle tyre makes and models from forensic-quality photographs of tyre impressions.

The project is designed as a **research prototype** to support forensic analysis workflows, particularly in the interpretation of tyre impression evidence in vehicle-related crime investigations.

---

## Academic Context

This project is submitted as part of an MSc dissertation in Forensic Investigation.

The research investigates:

> The application of computer vision and pattern recognition techniques for identifying vehicle tyres from impression evidence.

TyreMatch is not intended as a production forensic system but as an **experimental platform** for testing database-backed image analysis workflows and future machine learning integration.

---

## Features

- Flask-based REST API backend
- MySQL 8.0 relational database integration
- Database migrations via Flask CLI / Flask-Migrate
- Fully containerised Docker architecture
- Isolated test environment using ephemeral database instances
- Pytest-based automated testing framework
- Designed for future computer vision / ML pipeline integration

---

## System Architecture

The system is split into two Docker Compose environments:

### Production / Development Stack

- Flask backend service (`tyre_match_app`)
- MySQL database (`tyre_match_mysql`)
- Persistent Docker volume for database storage
- Automatic database migrations on startup
- Live code mounting for development

### Testing Stack

- Dedicated test runner container (`tyre_match_test`)
- Ephemeral MySQL test database (`tyre_match_test_mysql`)
- `tmpfs` in-memory storage for fast reset
- Fully isolated from development/production data

---

## Environment Configuration

The project uses a `.env` file for configuration. A template is provided in `.env.dist`.

### Example `.env.dist`

```env
# Environment to run in. Options [production, development]
ENVIRONMENT=development

# Application host configuration
HOST=0.0.0.0
PORT=7788

# Database configuration
DB_USER=local_user
DB_PASSWORD=password
DB_DRIVER=mysql
DB_NAME=tyre_match
DB_HOST=tyre_match_mysql
DB_PORT=3306

# External access ports (host machine mapping)
TEST_DB_HOST=0.0.0.0
EXPOSE_PORT=7788
EXPOSE_DB_PORT=56740
TEST_EXPOSE_DB_PORT=56741

#Parameters for celery asynchronous worker
CELERY_EXPOSE_PORT=6379
CELERY_PORT=6379
CELERY_BROKER_URL=redis://tyre_match_redis:6379/0
CELERY_RESULT_BACKEND=redis://tyre_match_redis:6379/0

# Flask configuration
FLASK_APP=main.py
```

---

## Notes
- EXPOSE_PORT is mapped to the Flask container port for host access
- EXPOSE_DB_PORT exposes MySQL for local debugging if required
- TEST_EXPOSE_DB_PORT isolates test database from development database
- DB_HOST uses Docker service name for internal networking
- HOST=0.0.0.0 ensures the Flask app is accessible inside Docker

--- 

## Requirements
- Docker
- Docker Compose v2+
- GNU Make (recommended)

--- 

## Running the Application
### Build and start the system
```bash
    make build
```

This will:

- Build the Flask backend container
- Start the MySQL database
- Run database migrations automatically
- Launch the API service

The API will be available at:
http://localhost:${EXPOSE_PORT}

--- 

## Running Tests
### Run test suite
```bash
    make test
```

This will:

- Spin up an isolated MySQL test database
- Run database migrations
- Execute the full pytest suite
- Tear down containers after execution 

### Run tests with coverage
```bash
    make test-coverage
```

Coverage report is generated at:

test-artifacts/coverage/index.html

---

## Docker Services
### Backend (tyre_match_app)
- Built from docker/api/Dockerfile
- Flask entrypoint: main:create_app
- Runs migrations before startup
- Mounted volume for development

### Database (tyre_match_mysql)
- MySQL 8.0
- Persistent Docker volume
- Health-checked startup dependency

--- 

## Testing Architecture
### Test Runner (tyre_match_test)
- Executes pytest inside container
- Runs migrations before tests
- Environment flag: TESTING=1

### Test Database (tyre_match_test_mysql)
- MySQL 8.0 ephemeral instance
- Uses tmpfs for fast in-memory reset
- Fully isolated per test run

--- 

## Makefile Commands
| Command            | Description                                 | 
|--------------------|---------------------------------------------|
| make build         | Build and start full application stack      |
| make test          | Run full test suite in isolated environment |
| make test-coverage | Run tests with coverage reporting           |

---

## Database Migrations

Migrations run automatically on startup:

```bash
  flask --app main:create_app db upgrade
```

Manual execution:
```bash
  docker compose exec tyre_match_app flask db upgrade
```

--- 

## Project Structure
```
.
├── .github/
|   └── workflows/
├── api/
|   |── repsponses/
|   └── routes/
├── celery_config/
└──database
│   |── migrations/
│   |── models/
│   └── repositories/
├── docker/
│   └── api/
|── domain/
|── files/
|    |── test_directory/    # Generated by test suite
|    └── tyre_impressions/  # Generated by tyre impression service
|── logging_config/
├── pipelines/
|   └── processors/
|── policies/
|── preprocessing/
|    └── processors/
|── services/
|── tasks/
|── test-artifacts/
├── tests/
|    |── api/
|    |   └── routes/
|    |── helpers/
|    |   └── factories/
|    |── mocks/
|    |   |── data/
|    |   |── database
|    |   |   └── repositories/
|    |   └── services/
|    |── models/
|    ├── pipelines/
|    |   └── processors/
|    |── policies/
|    |── repositories/
|    |── services/
|    └── utils/
├── utils/
├── .coverage     # Generated by test coverage
├── .coveragerc   # Generated by test coverage
├── .dockerignore
├── .env.dist
├── .gitignore
├── CHANGELOG.md
├── config.py
├── docker-compose.yml
├── docker-compose-test.yml
├── main.py
├── Makefile
├── README.md
└── requirements.txt
```

--- 

## Research Notes
- This is an MSc forensic investigation research prototype
- Intended for experimental evaluation of tyre pattern recognition workflows
- Not validated for operational forensic deployment
- Accuracy depends heavily on image quality and dataset completeness
- Designed as a foundation for:
  - Computer vision feature extraction
  - Pattern matching and similarity scoring
  - Future machine learning classification models

--- 

## License

This project is developed for academic purposes as part of an MSc dissertation. Licensing is subject to institutional requirements and may be updated after submission.

---