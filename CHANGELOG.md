## v0.1.4 `DATE`
### [Test Coverage]
- Statements: 
- Missed:  
- Coverage: %

### [Changed]
- Updated README project structure

## v0.1.3 27/05/2026
### [Test Coverage]
- Statements: 768
- Missed: 7
- Coverage: 99%

Reason for missed statements:
- New missing statements are in celery worker config and not tested due to this being for asynchronous tasks
- New missing statements are also in base processor, they are method declarations used for typing and have no content

### [Added]
- Added new get by id method for tyre impression processing to fetch files
- Created file factory for tests
- Added normalisation processor using Contrast Limited Adaptive Histogram Equalisation (CLAHE)
- Setup preprocessing pipeline in asynchronous task
- Added openCV dependencies to dockerfile

### [Changed]
- Override tyre impression repo get_by_id to include join for processing
- Changed logger from using flask app context to generic logger so it can be used by async tasks

### [Fixed]
- Fixed celery configuration and moved to `celery_config` directory

### [Removed]
- Removed explicit relationship definition between tyre impression processing and files
- Removed all instances of global importing with __init__.py to prevent packages being imported before they are needed


## v0.1.2 26/05/2026
### [Test Coverage]
- Statements: 635
- Missed: 4
- Coverage: 99%

### [Added]
- Created files table

### [Changed]
- Updated tyre impression model to reflect new relationship
- Updated tyre impression processing model to reflect files changes

### [Removed]
- Removed the tyre impression task service
- Removed tests/mocks to be more consistent with using patch methods

## v0.1.1 26/05/2026
### [Test Coverage]
- Statements: 553
- Missed: 5
- Coverage: 99%

Reason for missed statement:
- New missing statements are invocations of methods that do not exist yet or have no function 

### [Added]
- Created tyre impression processing table
- Created tyre impression processing service and repository
- Added celery configuration for asynchronous task processing
- Added tests for tyre impression processing repository

### [Changed]
- Updated project structure in readme
- Updated values in tyre impression status enum
- Updated tyre impression and tyre model tests to use factories
- Omitted celery_app.py from test coverage report

### [Fixed]
- Fixed tyre model service delete method to correctly use unit of work


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