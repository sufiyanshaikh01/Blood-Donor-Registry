# Blood-Donor-Registry
Blood Donor Registry <br>
A full-stack web application designed to manage blood donor registrations and search for donors by blood group. This project was built to practice API development and frontend-backend integration.

🛠️ Tech Stack <br>
Frontend: React.js

Backend: Go (Golang)

Database: CSV File (donors.csv)

Version Control: Git & GitHub

📂 Project Structure <br>
backend/: Contains the Go source code (main.go) and the CSV file used for data storage.

forntend/: Contains the React application, including the main logic in App.js.

🚀 How to Run <br>
1. Backend (Go) <br>

Navigate to the backend directory: cd backend.
Run the server: go run main.go.
The API will be available at http://localhost:8080.

**2. Frontend (React)** <br>

Navigate to the frontend directory: cd forntend.<br>
Install dependencies: npm install.<br>
Start the development server: npm start.<br>

🔌 API Features <br>
Fetch Donors: Retrieves a list of donors from the Go server.

Filter: Allows filtering donors based on specific blood groups using query parameters.
