## v0.2.1 `DATE`

### [Added]
- Added tyre models to the image processing pipeline

### [Removed]
- Removed redundant category, pattern type, and vehicle tpe fields from tyre model

## v0.2.1 30/08/2026

### [Added]
- Added binary processor to image processing pipeline, using adaptive thresholding to convert grayscale image to black/white background/foreground
- Added source image validation to base processor
- Added colour image constructor to test helpers
- Added assert result method to cv test case struct
- Added normalisation processor to image processing pipeline, using
``` normalised = (source / illumination) * mean(illumination) ``` to correct for broad lighting changes
- Added region of interest (ROI) isolation to the normalisation processor, hard-coded bounding box parameter and aggressive gaussian blur to remove detail outside ROI
- Added ROI parameters to the create method for tyre impressions

### [Changed]
- Updated example.jpg and example.jpeg to more accurately reflect the kind of image this system will receive
- Reordered steps in enhancement processor to reduce high-frequency noise from background
- Updated GitHub workflow to successfully build project with GoCV/OpenCV and run all tests on push to main branch

## v0.2.0 28/08/2026

### [Changed]
- Converted project from Python to Go to make better use of my familiarity with the language

### [Removed]
- Removed the normalisation processor from the pipeline
- Removed future test coverage reports from CHANGELOG.md, test coverage will now be reviewed manually and no longer reported

## v0.1.6 11/07/2026
### [Test Coverage]
- Statements: 944
- Missed: 5
- Coverage: 99%

### [Changed]
- Moved tyre impression metadata and file relationships to tyre impression model

### [Fixed]
- Fixed tyre model upload tests creating files in the wrong directory

### [Removed]
- Removed TyreImpressionProcessing object and associated artifacts

## v0.1.5 07/06/2026
### [Test Coverage]
- Statements: 1011 
- Missed: 5
- Coverage: 99%

### [Added]
- Added get by ID method for tyre impressions
- Added file route for returning images to the frontend
- Added image upload endpoint for tyre models

### [Fixed]
- Stopped the file service from appending the file name to the end of the returned file location
- Fixed some file naming inconsistencies


## v0.1.4 02/06/2026
### [Test Coverage]
- Statements: 876
- Missed: 5
- Coverage: 99%

### [Added]
- Added normalisation processor
- Added resize method to normalisation processor that sets image to specified target width/height while maintaining aspect ratio
- Added correct skew method to normalisation processor to rotate image relative to strong lines detected by OTSU masking and calculation a principal axis using eigen values
- Added enhancement processor and moved CLAHE logic to there
- Added more extensive tests for enhancement processor, including deterministic assertions
- Added denoise and sharpen methods to the enhancement processor

### [Changed]
- Updated README project structure
- Moved process and transform processor methods to base processor since they will need to be on all processors
- Updated base processor transform to allow for multiple stages of transformation
- Updated return types in processing pipeline to reduce IO Read/Write operations
- Changed the tyre impression processed image extensions from jpg to png to reduce cumulative compression artifacts

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