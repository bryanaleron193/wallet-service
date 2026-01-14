# wallet-service

Features:
1. User Authentication: JWT-based login system.
2. Balance Inquiry API: Balance checking.
3. Withdraw Transaction: Withdrawal operations.
4. Observability: Health check endpoint.
5. Documentation: Interactive Swagger UI.

Tech Stack:
1. Language: Go 1.25
2. Framework: Echo v4
3. Database: PostgreSQL 18 with pgx/v5
4. Documentation: Swaggo
5. Live Reload: Air
6. Containerization: Docker & Docker Compose

Step:
1. docker-compose up --build
2. make migrate-up (do mingw32-make migrate-up if you use windows)
3. make seed (do mingw32-make seed if you use windows)

Installation & Setup
1. Clone the repository:
    - git clone https://github.com/bryanaleron193/wallet-service.git
2. Start the Database + Application:
    - docker-compose up --build
3. Run Database Migrations:
    - make migrate-up
    # Windows: mingw32-make migrate-up
4. Seed Test Data: This creates a test user with an initial balance.
    - make seed
    # Windows: mingw32-make seed

API Documentation:
1. The interactive API documentation is pre-generated and available via Swagger UI. You do not need to install any additional tools to view it.
2. URL: http://localhost:8081/swagger/index.html
3. Key Endpoints:
- GET /health: Check service status.
- POST /login: Get your JWT Bearer token.
- GET /api/v1/wallets/balance: Check current balance (Requires Auth).
- POST /api/v1/wallets/withdraw: Withdraw funds (Requires Auth).

Testing the API
1. Login: Use the seeded username to obtain a JWT token.
2. Authorize: Click the "Authorize" button in Swagger and enter Bearer <your_token>.
3. Transaction: Perform a withdrawal and balance inquiry.