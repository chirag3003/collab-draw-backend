# Collab Draw Backend

A robust GraphQL API built with Go for real-time collaborative drawing applications. This backend powers the `collab-draw` platform, handling workspace management, project persistence, and live synchronization using GraphQL subscriptions.

## 🚀 Tech Stack

- **Language:** Go 1.25+
- **GraphQL Engine:** [gqlgen](https://gqlgen.com/)
- **Database:** MongoDB (via official mongo-driver/v2)
- **Authentication:** [Clerk](https://clerk.com/) (Clerk SDK for Go)
- **Real-time:** Gorilla WebSockets for GraphQL Subscriptions
- **Router:** Chi

## 📂 Project Structure

- `graph/`: GraphQL schema definitions (`.graphqls`) and generated code.
  - `resolvers/`: Implementation of GraphQL queries, mutations, and subscriptions.
- `internal/`: Core business logic and infrastructure.
  - `auth/`: Clerk middleware for session verification.
  - `db/`: MongoDB connection management.
  - `models/`: BSON-tagged Go structs for database entities.
  - `repository/`: Data access layer (Repository Pattern) for Projects and Workspaces.
- `server.go`: Application entry point and server configuration.

## 🛠️ Getting Started

### Prerequisites

- Go 1.25 or higher
- MongoDB instance (local or Atlas)
- Clerk account and API keys

### Environment Variables

Create a `.env` file in the root directory:

```env
PORT=8080
MONGO_URI=mongodb://localhost:27017
DATABASE_NAME=collab-draw
CLERK_SECRET_KEY=sk_test_...
```

### Installation

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Run the server:
   ```bash
   go run server.go
   ```

The GraphQL playground will be available at `http://localhost:8080/`.

## 📡 Real-time Collaboration

The backend uses GraphQL Subscriptions to sync drawing elements across clients. 
- **Mutation:** `updateProject` updates the drawing state in MongoDB and broadcasts to subscribers.
- **Subscription:** `project(id: ID!)` allows clients to receive live updates.

For more details, see [SUBSCRIPTION_GUIDE.md](./SUBSCRIPTION_GUIDE.md).

## 🔒 Authentication

All requests (except the playground in dev mode) require a valid Clerk session token passed in the `Authorization` header:
`Authorization: Bearer <your_session_token>`
