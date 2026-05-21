## v0.0.5 `DATE`
### [Test Coverage]
- Statements: 
- Missed: 
- Coverage: % 

### [Added]
- Created a tyre impression model with a status enum data type
- Created database migration for new `tyre_impressions` table
- Added a file service with file naming policies
- Added get_all endpoint for tyre impressions
- Added image upload endpoint for tyre impressions

### [Changed]
- Changed styling on subheadings in changelog
- Updated endpoints for consistency on plural/singular words

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