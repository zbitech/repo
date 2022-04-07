# ZBI Repository Service
The repository service is a go-lang package the implements the persistence and iam services for the 
Zcash Blockchain Infrastructure (ZBI) project. It has dependencies on the common package.

### Dependencies
- http://github.com/zbitech/common
- http://github.com/mongodb/mongo-go-driver

## Persistence Mechanism
ZBI requires a persistence mechanism to store its users, meta-data about its resources
and the relationship between users and those resources. A repository service interface 
in the common package defines the functionality that is expected. This package provides
two default implementations as described below. Users can substitute with implementations
of their choice.

### Memory-Based Persistence
The memory-based service is intended for light-weight projects and especially environments
with limited resources.

### Database Persistence
The database service is intended for production-like environments. This is a default
implementation that uses MongoDB NoSQL database.

## Identity & Access Management
ZBI requires users to provide an authentication token (JWT) or API key when accessing endpoints.
It uses an IAM service to manage authentication services for its users.

### Basic Authentication
The Basic authentication service is intended for light-weight projects that do not require
an OAuth or OIDC server. This is a default implementation that stores hashed credentials
in the data store and requires users to authenticate using the Basic authentication scheme.
It also manages API keys in the data store.

### OIDC/OAuth2 Authentication
This service is intended for production-like environments 

## Access Authorizer

## Contributing
Pull requests are welcome. for major changes, please open an issue first to discuss what you would like to change.
## License
[MIT](https://choosealicense.com/licenses/mit/)