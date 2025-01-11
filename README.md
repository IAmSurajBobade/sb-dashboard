# Go Service App

This project is a simple Go web service that allows users to create and retrieve entries associated with a unique UUID and a 4-digit password. Users can input a list of dates and a name, which will be stored and can later be accessed via a GET request.

## Project Structure

```
go-service-app
├── src
│   ├── main.go          # Entry point of the application
│   ├── handlers
│   │   ├── create.go    # Handler for the /create endpoint
│   │   └── get.go       # Handler for the /get endpoint
│   ├── models
│   │   └── data.go      # Data structures for the application
│   ├── routes
│   │   └── routes.go     # Routing setup for the application
│   └── templates
│       ├── create.html   # HTML template for creating entries
│       └── get.html      # HTML template for retrieving entries
├── go.mod                # Module definition and dependencies
└── README.md             # Project documentation
```

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd go-service-app
   ```

2. **Install dependencies:**
   ```
   go mod tidy
   ```

3. **Run the application:**
   ```
   go run src/main.go
   ```

4. **Access the service:**
   Open your web browser and navigate to `http://localhost:8080/create` to create a new entry or `http://localhost:8080/get` to retrieve an entry.

## Usage Examples

- **Creating an Entry:**
  - Navigate to the `/create` endpoint.
  - Fill out the form with your name, a list of dates, and a 4-digit password.
  - Submit the form to create a new entry.

- **Retrieving an Entry:**
  - Navigate to the `/get` endpoint.
  - Enter the UUID and password associated with your entry.
  - Submit the form to view the number of days since the specified dates.

## Technologies Used

- Go (Golang)
- HTMX for dynamic HTML interactions

## License

This project is licensed under the MIT License. See the LICENSE file for more details.