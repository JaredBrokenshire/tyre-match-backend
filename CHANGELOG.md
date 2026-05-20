## v0.0.3 `DATE` - Test Coverage: 
### Added
- Added pagination limit and offset to base repo get_all method
- Added `get_all` route for `tyre_model` and response objects
- Added `get_by_id` route for `tyre_model`
### Changed
- Fixed `.env.dist` to match required environment variables
- Omitted `database/session.py` and `database/extensions.py` from test coverage report

##  v.0.0.2 20/05/2026 - Test Coverage: 97%

### Added
- Added Pytest functionality (with coverage reporting) using separate docker-compose for mock database
- Added generic base repository format with CRUD functions
- Added tyre model repo inheriting from base repo
- Added tyre model and subsequent database migration to create `tyre_models` table

### Changed
- Updated README to include project details and makefile commands

## v0.0.1 - 18/05/2026

### Added
- Project initialised with flask server connected to MySQL database using docker-compose
- Flask-Migrate and SQLAlchemy initialised for DB representation and migrations