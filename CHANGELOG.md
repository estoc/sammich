# ChewCrew Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased][unreleased]
- Integration with a Places API such as Yelp or Google
- Ability to input addresses
- Finding places with the lowest average distance

## [0.0.1] - 2015-03-22
### Added
- Ready command, used to generate votes after everyone is Ready.
- API command, used as an api reference and functionality testing.
- CurrentVoter to the session object, used to help clients keep track of their own status.
- The server port flag.

### Changed
- Clearing out the Choices array, and also the Ready and Voted flags when appropriate to reduce the amount of data sent back.
- Only the 'Get' command returns a full session object, acting as the singular source of data.
- Errors are now returned as JSON, ex: {"error": "Session not found"}.

### Removed
- Start command, since it's now using the Ready model.

### Fixed
- Query strings not found now return an empty string instead of throwing an error.
