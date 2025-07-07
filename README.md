# Hotel Booking API

A full-featured hotel reservation REST API built with **Go**, **MongoDB**, and **Fiber**, supporting user authentication, hotel/room management, and booking operations.

---

## Features

* User registration and login with JWT authentication
* Role-based access (admin vs normal user)
* Secure password hashing using `bcrypt`
* Hotel and room management (add, fetch)
* Room booking and cancellation with overlap checks
* Unit test coverage for user and auth handlers
* Environment configuration via `.env`

---

## Technologies Used

* Go (Golang)
* MongoDB with `mongo-driver`
* Fiber framework
* JWT (via `golang-jwt`)
* Bcrypt for password hashing

---

## Environment Variables (`.env`)

```
MONGO_DB_URI=mongodb://localhost:27017
HTTP_LISTEN_ADDRESS=:3000
JWT_SECRET=somethingsupersecretthatNOBODYKNOWS
MONGO_DB_NAME=hotel-reservation
MONGO_TEST_DB_NAME=hotel-reservation-test
```

---

## Running the Project

```bash
# 1. Load environment variables

# 2. Run MongoDB locally (if not already running)
# e.g., via Docker: docker run -p 27017:27017 mongo

# 3. Run the project
make run
```

---

## API Endpoints

### Auth

* `POST /auth` — Login and get JWT

### Users

* `POST /users` — Register a new user
* `GET /users/:id` — Get user by ID
* `GET /users` — List all users (admin only)
* `PUT /users/:id` — Update user name
* `DELETE /users/:id` — Delete user

### Hotels

* `GET /hotels` — List all hotels
* `GET /hotels/:id` — Get hotel by ID
* `GET /hotels/:id/rooms` — List rooms in hotel

### Rooms

* `GET /rooms` — List all rooms
* `POST /rooms/:id/book` — Book a room (requires authentication)

### Bookings

* `GET /bookings` — List all bookings (admin only)
* `GET /bookings/:id` — Get booking details (owner or admin only)
* `POST /bookings/:id/cancel` — Cancel a booking (owner or admin only)

---

## Core File Structure

* `api/` — Fiber route handlers for users, hotels, rooms, and bookings
* `db/` — MongoDB storage interfaces and implementations
* `types/` — Core domain models and validation
* `fixtures/` — Helpers to seed test or demo data

---

## Testing

Unit tests included for:

* User creation
* Authentication

Run tests:

```bash
make test
```

---

## Example Usage

Authenticate with already created user:

```json
POST /auth
{
  "email": "qwe@rty.com",
  "password": "qwe_rty"
}
// returns { user, token }
```

Use `X-Api-Token: <JWT>` header in subsequent requests.

---

## Author

Created by [AyanokojiKiyotaka8](https://github.com/AyanokojiKiyotaka8) — open for feedback and collaboration!

---

## License

MIT License — free to use, distribute, and modify with attribution.
