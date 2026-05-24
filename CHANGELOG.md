## v0.1.1 `DATE`
### [Test Coverage]
- Statements: 
- Missed: 
- Coverage: %



## v0.1.0 24/05/2026
### [Test Coverage]
- Statements: 417
- Missed: 1
- Coverage: 99%

Reason for decrease: Missing statement is an if statement used as a flag for the testing environment in /database/unit_of_work.py 

### [Added]
- Added a mock base repository to prevent service tests from persisting data in the database
- Added update and delete endpoints for tyre models
- Added model not found error exception
- Added assertion helper methods for tests
- Added model factories for tests

### [Changed]
- Updated tyre impression upload route and tyre impression service to give more detailed logs and error outputs
- Changed the working directory in docker-compose
- Moved file service to a class object
- Updated `dataset_source` and `notes` columns in `tyre_models` to be longtext
- Updated lengths of varchar fields in `tyre_models` table
- Refactored all tests to help with future debugging and maintainability
- Moved search parameter out of base repo to keep it generic

### [Removed]
- Removed debug logs from tyre model service


## v0.0.5 22/05/2026
### [Test Coverage]
- Statements: 288
- Missed: 0
- Coverage: 100% 

### [Added]
- Created a tyre impression model with a status enum data type
- Created database migration for new `tyre_impressions` table
- Added a file service with file naming policies
- Added get_all endpoint for tyre impressions
- Added image upload endpoint for tyre impressions
- Added custom exceptions
- Added `TyreImpressionService` to handle upload image logic
- Added `TyreModelService` to handle tyre model endpoint logic
- Added create endpoint for tyre models
- Added tyre model service for create logic

### [Changed]
- Changed styling on subheadings in changelog
- Updated endpoints for consistency on plural/singular words
- Updated `.dockerignore` to include `test-artifacts` and coverage reports

## v0.0.4 21/05/2026
### [Test Coverage]
- Statements: 109
- Missed: 0
- Coverage: 100% 

### [Added]
- Added model_name or manufacturer search parameter to tyre_model list route
- Added appropriate CORS origins to allow requests from frontend

### [Changed]
- Changed pagination parameters to be page and page_size and used built-in paginate method

## v0.0.3 20/05/2026
### [Test Coverage]
- Statements: 100
- Missed: 0
- Coverage: 100%
### [Added]
- Added pagination limit and offset to base repo get_all method
- Added `get_all` route for `tyre_model` and response objects
- Added `get_by_id` route for `tyre_model`
### [Changed]
- Fixed `.env.dist` to match required environment variables
- Omitted `database/session.py` and `database/extensions.py` from test coverage report

##  v.0.0.2 20/05/2026 - Test Coverage: 97%

### [Added]
- Added Pytest functionality (with coverage reporting) using separate docker-compose for mock database
- Added generic base repository format with CRUD functions
- Added tyre model repo inheriting from base repo
- Added tyre model and subsequent database migration to create `tyre_models` table

### [Changed]
- Updated README to include project details and makefile commands

## v0.0.1 - 18/05/2026

### [Added]
- Project initialised with flask server connected to MySQL database using docker-compose
- Flask-Migrate and SQLAlchemy initialised for DB representation and migrations