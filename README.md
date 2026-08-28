# TyreMatch

TyreMatch is a computer vision backend system developed as part of an **MSc Forensic Investigation** dissertation research project.

The project investigates the feasibility of using computer vision and image-processing techniques to assist in identifying vehicle tyre makes and models from **forensic tyre impression photographs**.

The current implementation provides a RESTful backend for storing tyre models and tyre impressions, uploading impression photographs, and processing uploaded images through an OpenCV-based image enhancement pipeline. The project is intended to provide the foundation for subsequent feature extraction, image matching and similarity-scoring research.

> **Research prototype:** TyreMatch is an academic research project and has not been validated for operational forensic use. Results produced by the system must not be treated as independent evidence of forensic identification.

---

## Academic Context

This project is being developed as part of an **MSc dissertation in Forensic Investigation**.

The research investigates:

> **The application of computer vision and pattern recognition techniques for identifying vehicle tyres from impression evidence.**

The overall research objective is to investigate whether computer vision techniques can be used to compare an unknown tyre impression photographed at a crime scene against a known collection of tyre makes and models.

The system is currently being developed as an experimental platform rather than a production forensic application.

The project is intentionally structured so that additional computer vision processors can be introduced as the research progresses. Future stages are expected to investigate techniques such as feature extraction, feature matching and similarity scoring, including approaches based on OpenCV feature descriptors such as **ORB and SIFT**.

These future techniques are **not currently implemented** and are therefore outside the scope of the current application functionality.

---

## Dataset

The research dataset is **external to this repository** and is not committed to Git.

The dataset contains **96 tyre image sets**.

Each tyre set contains multiple representations of the same tyre, including photographs:

* Under different lighting conditions
* Of the physical tyre
* Of impressions produced on paper
* Of impressions produced in ground-simulating subtrates

The intention is to use these images as known reference data against which an uploaded unknown impression can eventually be compared.

Because the dataset is external, it must be obtained and managed separately from this repository. Any use of the dataset should comply with its applicable licensing, research and institutional requirements.

---

## System Architecture

TyreMatch uses a containerised backend architecture consisting of an application container and MySQL database services.

### Application Stack

The main application stack contains:

* **TyreMatch backend**

    * Go application
    * Echo HTTP server
    * GORM database layer
    * OpenCV/GoCV image processing
    * Local file storage

* **MySQL**

    * Persistent application database
    * Stores tyre models, tyre impressions and file metadata

### Testing Stack

A separate Docker Compose configuration provides a dedicated MySQL database for automated tests.

Tests create an in-process test server using Echo's HTTP handling mechanisms rather than requiring the API to be exposed as a separate running HTTP service.

This allows API behaviour, database operations and service interactions to be tested without modifying the development database.

---

## Technology Stack

| Component           | Technology                             |
| ------------------- | -------------------------------------- |
| Language            | Go 1.25                                |
| HTTP Framework      | Echo v4                                |
| ORM                 | GORM                                   |
| Database            | MySQL 5.7                              |
| Database Driver     | go-sql-driver/mysql                    |
| Migrations          | Project `go-migrations` implementation |
| Computer Vision     | OpenCV                                 |
| Go OpenCV bindings  | GoCV v0.43.0                           |
| OpenCV Docker image | `gocv/opencv:4.13.0`                   |
| API Documentation   | Swagger / OpenAPI                      |
| Testing             | Go standard `testing` package          |
| Assertions          | Testify                                |
| Containerisation    | Docker / Docker Compose                |

---

## Requirements

The recommended development environment requires:

* Docker
* Docker Compose
* Go 1.25 or later
* Git

For local execution of Go tests, Go must be installed on the host system.

For image-processing development outside Docker, a compatible native OpenCV installation is also required by GoCV.

The Docker development environment already provides OpenCV through the:

```text
gocv/opencv:4.13.0
```

base image.

---

## GoCV and OpenCV Versioning

TyreMatch uses:

```text
GoCV:   v0.43.0
OpenCV: 4.13.0
Go:     1.25
```

The GoCV version is defined in `go.mod`:

```text
gocv.io/x/gocv v0.43.0
```

The Docker image used to build and run the application is:

```text
gocv/opencv:4.13.0
```

These versions should be kept compatible when modifying the computer vision implementation.

In particular, changes to OpenCV or GoCV versions should be tested against the existing image-processing test suite because GoCV relies on the underlying native OpenCV installation.

---

## Environment Configuration

The application uses a `.env` file for runtime configuration.

A template is provided as:

```text
.env.dist
```

Create the local environment file before starting the application:

```bash
cp .env.dist .env
```

An example configuration is:

```env
# Application environment
ENVIRONMENT=development

# Application host configuration
HOST=localhost
PORT=7788

# Database configuration
DB_USER=local_user
DB_PASSWORD=password
DB_DRIVER=mysql
DB_NAME=tyre_match
DB_HOST=tyre-match-backend-mysql
DB_PORT=3306

# Host machine ports
TEST_DB_HOST=0.0.0.0
EXPOSE_PORT=7788
EXPOSE_DB_PORT=56742
TEST_EXPOSE_DB_PORT=56743
```

### Configuration notes

* `PORT` is the port used by the Go application inside the container.
* `EXPOSE_PORT` is the host port used to access the API.
* `DB_HOST` should use the Docker Compose MySQL service name when the application is running inside Docker.
* `DB_PORT` is the internal MySQL port.
* `EXPOSE_DB_PORT` exposes the development database to the host when direct database access is required.
* `TEST_EXPOSE_DB_PORT` exposes the test database to the host.
* Database credentials should be changed for environments where the application is not being used purely for local development.

Do not commit sensitive environment configuration to the repository.

---

## Running the Application

### 1. Create the environment file

```bash
cp .env.dist .env
```

Review the configuration before starting the application.

### 2. Build and start the application

```bash
docker-compose up --build -d
```

This builds the Go/OpenCV application image and starts the backend and MySQL database containers.

The backend will be available at:

```text
http://localhost:7788
```

assuming the default `EXPOSE_PORT` is being used.

The application container waits for the database before running the database migration process and starting the Go server.

---

## Running the Test Environment

Start the dedicated test database using:

```bash
docker-compose -f ./docker-compose-test.yml up --build -d
```

The test suite can then be executed with:

```bash
go test ./tests
```

The tests use the dedicated test database rather than the development database.

### Test architecture

The tests use:

* A dedicated MySQL test database
* An in-process Echo HTTP server
* Mock file storage
* Mock image-processing dependencies where appropriate
* Test factories for database models
* Test helpers for HTTP requests and assertions

Tests are contained within:

```text
./tests
```

---

## Testing

The project currently uses the standard Go testing framework.

Run all tests with:

```bash
go test ./tests
```

Tests cover areas including:

* Tyre impression API behaviour
* Tyre model API behaviour
* File API behaviour
* Service behaviour
* Validation
* Image-processing behaviour
* Image enhancement processors
* Database interactions
* Error handling
* File storage interactions

The image-processing tests also verify properties of the processing operations, including that:

* Processing does not unexpectedly mutate the source image
* Processing is deterministic
* CLAHE modifies local image contrast
* Sharpening modifies image characteristics
* Empty and invalid images are handled appropriately

Test coverage reports have previously been generated using Go's built-in coverage functionality. Coverage reporting is not currently treated as a primary project metric; test coverage is being assessed manually during development.

---

## API

TyreMatch exposes a RESTful HTTP API.

The API is divided primarily into:

* Tyre models
* Tyre impressions
* Files

### Tyre Models

```text
GET    /tyre-models
POST   /tyre-models
GET    /tyre-models/:id
PUT    /tyre-models/:id
DELETE /tyre-models/:id
POST   /tyre-models/:id
```

Tyre model records can contain information including:

* Manufacturer
* Model name
* Category
* Vehicle type
* Width
* Aspect ratio
* Rim diameter
* Groove count
* Pattern type

The final `POST /tyre-models/:id` route is used for uploading files associated with a tyre model.

### Tyre Impressions

```text
GET    /tyre-impressions
POST   /tyre-impressions
GET    /tyre-impressions/:id
DELETE /tyre-impressions/:id
POST   /tyre-impressions/:id
```

A tyre impression contains information including:

* Pixels per inch
* Processing status
* Edge density
* Void ratio
* Groove count
* Associated image files

The impression upload endpoint accepts a multipart file upload using the field:

```text
file
```

### Supported image formats

The current upload validator accepts:

```text
JPEG
JPG
PNG
WEBP
```

Both the file extension and detected MIME type are validated.

PDF and other non-image files are rejected.

---

## File Storage

Uploaded files are stored using the project's local filesystem storage implementation.

Tyre impression files are organised using paths based on the impression ID and file type.

For example:

```text
tyre-impressions/
└── <impression-id>/
    ├── original/
    └── enhanced/
```

File metadata is stored in the MySQL `files` table.

Generated processing files are encoded as PNG images before being stored.

---

## API Documentation

Swagger/OpenAPI documentation is maintained in:

```text
docs/
├── swagger.yaml
├── swagger.json
└── docs.go
```

The Docker build process regenerates the Swagger documentation using the annotations contained in the Go source code.

The Swagger specification can therefore be used as the reference for the currently exposed API schema.

---

## Project Structure

```text
.
├── .github/
│   └── workflows/
│       └── ci.yml
│
├── api/
│   ├── handlers/
│   │   ├── file_handler.go
│   │   ├── tyre_impression_handler.go
│   │   └── tyre_model_handler.go
│   ├── requests/
│   │   ├── file_requests.go
│   │   ├── tyre_impression_requests.go
│   │   └── tyre_model_requests.go
│   ├── responses/
│   │   ├── file_responses.go
│   │   ├── response_wrapper.go
│   │   ├── tyre_impression_responses.go
│   │   └── tyre_model_responses.go
│   ├── routes/
│   │   ├── file.go
│   │   ├── routes.go
│   │   ├── tyre_impression.go
│   │   └── tyre_model.go
│   └── server.go
│
├── assets/
│   └── # Test/example files
│
├── cv/
│   └── processors/
│       ├── base_processor.go
│       └── enhancement_processor.go
│
├── db/
│   ├── migrations/
│   │   ├── entry.go
│   │   ├── list/
│   │   └── process/
│   ├── models/
│   │   ├── file.go
│   │   ├── tyre_impression.go
│   │   └── tyre_model.go
│   ├── repositories/
│   │   ├── file_repository.go
│   │   ├── repository.go
│   │   ├── repos.go
│   │   ├── tyre_impression_repository.go
│   │   └── tyre_model_repository.go
│   └── db.go
│
├── docker/
│   └── api/
│       └── Dockerfile
│
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── pkg/
│   ├── closer/
│   ├── dependencies/
│   ├── file_storage/
│   ├── go-migrations/
│   ├── try/
│   └── validation/
│
├── services/
│   ├── errors.go
│   ├── file_service.go
│   ├── image_processing_service.go
│   ├── services.go
│   ├── tyre_impression_service.go
│   ├── tyre_model_service.go
│   └── upload_service.go
│
├── tests/
│   ├── cv_tests/
│   ├── factories/
│   ├── helpers/
│   ├── mocks/
│   ├── service_tests/
│   ├── base_test.go
│   ├── file_test.go
│   ├── tyre_impression_test.go
│   └── tyre_model_test.go
│
├── .dockerignore
├── .env.dist
├── .gitignore
├── CHANGELOG.md
├── docker-compose-test.yml
├── docker-compose.yml
├── go.mod
├── go.sum
├── VERSION
├── version.go
└── README.md
```

---

## Development Workflow

A typical development workflow is:

### Start the application

```bash
docker-compose up --build -d
```

### Start the test database

```bash
docker-compose -f ./docker-compose-test.yml up --build -d
```

### Run tests

```bash
go test ./tests
```

---

## Continuous Integration

The project contains a GitHub Actions workflow:

```text
.github/workflows/ci.yml
```

The CI workflow currently:

1. Checks out the repository
2. Creates the `.env` file from `.env.dist`
3. Sets up Docker Buildx
4. Builds the application and test containers
5. Runs the Go test suite

The CI pipeline is intended to provide a basic automated check that the application continues to build and the test suite passes.

---

## Research Limitations

The current implementation has several important limitations.

### Dataset dependency

The research dataset contains 96 tyre sets and is external to the application.

The effectiveness of any eventual matching approach will therefore depend on:

* Dataset size
* Dataset diversity
* Image quality
* Image acquisition conditions
* Lighting conditions
* Impression substrate
* Tyre wear
* Camera characteristics
* Image scale and resolution
* Similarity between the unknown tyre and reference images

### No forensic validation

The system has not undergone forensic validation and does not establish the evidential reliability required for operational forensic deployment.

Any experimental matching results should therefore be interpreted as research outputs rather than forensic conclusions.

---

## Important Forensic Disclaimer

TyreMatch is an **experimental academic research prototype**.

It has not been validated for use as an operational forensic identification system and should not be used to make independent forensic identification decisions.

Computer vision similarity does not, by itself, establish that a crime-scene impression originated from a particular tyre.

Any future matching results produced by TyreMatch should be treated as **candidate-generation or research results** requiring appropriate forensic interpretation, validation and corroboration.

The system's eventual performance will depend on the quality of the underlying dataset, image acquisition conditions, preprocessing methods, feature extraction techniques and matching methodology.

---

## License

TyreMatch is developed for academic purposes as part of an **MSc Forensic Investigation dissertation**.

Licensing is subject to institutional requirements and may be updated following completion of the research project.
