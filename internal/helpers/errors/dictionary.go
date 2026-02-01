package errors

var (
	ERROR_INTERNAL_SERVER      = New("Internal server error occurred")
	ERROR_UNPROCESSABLE_ENTITY = New("The request was well-formed but was unable to be followed due to semantic errors")
	ERROR_BAD_REQUEST          = New("The server could not understand the request due to invalid parameters")
	ERROR_NOT_FOUND            = New("The requested resource was not found")
	ERROR_UNAUTHORIZED         = New("The request requires user authentication")

	ERROR_MAIN_COMMAND                  = New("Whoops. There was an error while executing your CLI")
	ERROR_CONNECTION                    = New("Unable to connect to the database")
	ERROR_REQUIRED_USERNAME_EMAIL_INPUT = New("Username or Email is required")
	ERROR_INVALID_PASSWORD              = New("The password does not meet the required criteria")
	ERROR_LOGIN_FAILED                  = New("Login failed. Please check your credentials and try again.")
	ERROR_CLIENT_ID_NOT_FOUND           = New("Client ID not found in context, please provide a valid Client ID")
	ERROR_CREATE_ACCESS_TOKEN           = New("Failed to create access token")
	ERROR_CREATE_SESSION_TOKEN          = New("Failed to create session token")
	ERROR_DUPLICATED_KEY                = New("Duplicate key value violates unique constraint")

	ERROR_MISSING_CLIENT_ID      = New("Missing Client ID in the request")
	ERROR_MISSING_USER_ID        = New("Missing User ID in the request")
	ERROR_MISSING_APPLICATION_ID = New("Missing Application ID in the request")
	ERROR_MISSING_TOKEN_ID       = New("Missing Token ID in the request")
)

var (
	BAD_REQUEST = []error{
		ERROR_REQUIRED_USERNAME_EMAIL_INPUT,
		ERROR_INVALID_PASSWORD,
		ERROR_LOGIN_FAILED,
		ERROR_BAD_REQUEST,
	}

	NOT_FOUND = []error{
		ERROR_NOT_FOUND,
	}

	UNPROCESSABLE_ENTITY = []error{
		ERROR_UNPROCESSABLE_ENTITY,
		ERROR_DUPLICATED_KEY,
	}

	INTERNAL_SERVER = []error{
		ERROR_CONNECTION,
		ERROR_MAIN_COMMAND,
		ERROR_CREATE_SESSION_TOKEN,
		ERROR_CREATE_ACCESS_TOKEN,
		ERROR_INTERNAL_SERVER,
	}

	UNAUTHORIZED = []error{
		ERROR_CLIENT_ID_NOT_FOUND,
		ERROR_MISSING_CLIENT_ID,
		ERROR_MISSING_USER_ID,
		ERROR_MISSING_APPLICATION_ID,
		ERROR_MISSING_TOKEN_ID,
		ERROR_UNAUTHORIZED,
	}
)
