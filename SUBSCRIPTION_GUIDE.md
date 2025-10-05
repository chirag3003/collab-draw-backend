# Real-time Project Subscriptions

## Overview
The backend now supports real-time GraphQL subscriptions for project updates. When any user mutates a project (creates or updates), all users subscribed to that project will receive the updated project data in real-time via WebSocket.

## Implementation Details

### Components Added/Modified:

1. **Resolver with Subscription Management** (`graph/resolvers/resolver.go`)
   - Added `projectSubscribers` map to track active subscriptions per project
   - Added `subscribersMutex` for thread-safe subscription management
   - Implemented `subscribeToProject()` to add new subscribers
   - Implemented `unsubscribeFromProject()` to remove subscribers and clean up
   - Implemented `broadcastProjectUpdate()` to send updates to all subscribers

2. **Project Resolvers** (`graph/resolvers/project.resolvers.go`)
   - Modified `CreateProject` mutation to broadcast newly created projects
   - Modified `UpdateProject` mutation to broadcast project updates
   - Implemented `Project` subscription resolver with:
     - Authentication check to verify user access
     - Initial project state delivery
     - Automatic cleanup on connection close

3. **Server Configuration** (`server.go`)
   - Updated to use `NewResolver()` function for proper initialization
   - WebSocket transport already configured with keep-alive pings

## How to Use

### Client-Side Subscription Example

#### Using Apollo Client (React/JavaScript)
```javascript
import { gql, useSubscription } from '@apollo/client';

const PROJECT_SUBSCRIPTION = gql`
  subscription OnProjectUpdate($id: ID!) {
    project(id: $id) {
      id
      name
      description
      owner
      workspace
      personal
      appState
      elements
      createdAt
    }
  }
`;

function ProjectEditor({ projectId }) {
  const { data, loading, error } = useSubscription(
    PROJECT_SUBSCRIPTION,
    { variables: { id: projectId } }
  );

  if (loading) return <p>Connecting...</p>;
  if (error) return <p>Error: {error.message}</p>;

  const project = data?.project;
  // Render your project editor with live updates
  return <div>{/* Your project UI */}</div>;
}
```

#### WebSocket Connection Setup
```javascript
import { GraphQLWsLink } from '@apollo/client/link/subscriptions';
import { createClient } from 'graphql-ws';

const wsLink = new GraphQLWsLink(
  createClient({
    url: 'ws://localhost:8080/query',
    connectionParams: {
      // Include authentication headers
      authorization: `Bearer ${yourAuthToken}`,
    },
  })
);
```

### Testing with GraphQL Playground

1. Open the GraphQL Playground at `http://localhost:8080/`
2. Add authentication header in the bottom left:
   ```json
   {
     "Authorization": "Bearer your-clerk-token"
   }
   ```
3. Create a subscription:
   ```graphql
   subscription {
     project(id: "your-project-id") {
       id
       name
       appState
       elements
     }
   }
   ```
4. In another tab, mutate the project:
   ```graphql
   mutation {
     updateProject(
       id: "your-project-id"
       appState: "new state"
       elements: "new elements"
     )
   }
   ```
5. You'll see the subscription tab receive the update automatically!

## Features

✅ **Real-time Updates**: All subscribers receive updates instantly when projects change
✅ **Authentication**: Subscriptions verify user access before allowing connections
✅ **Initial State**: Subscribers receive the current project state immediately upon connecting
✅ **Automatic Cleanup**: Subscriptions are properly cleaned up when clients disconnect
✅ **Thread-Safe**: Uses mutex locks to safely handle concurrent subscriptions
✅ **Efficient Broadcasting**: Only sends updates to subscribers of the specific project
✅ **Non-blocking**: Uses Go channels with select statements to prevent blocking

## Architecture

```
Client 1 ──┐
           ├──> Subscribe(projectID) ──> Channel 1 ──┐
Client 2 ──┤                                          ├──> Resolver.projectSubscribers[projectID]
           └──> Subscribe(projectID) ──> Channel 2 ──┘

Mutation (UpdateProject) ──> broadcastProjectUpdate() ──> All Channels receive update
```

## Security

- All subscriptions require authentication via Clerk JWT
- Users can only subscribe to projects they have access to
- Access is verified on subscription creation
- WebSocket connections use the same CORS policy as HTTP requests

## Performance Considerations

- Channels are buffered (size 1) to prevent blocking
- Failed channel sends (full/closed channels) are gracefully skipped
- Empty subscriber lists are automatically cleaned up
- Mutex locks are minimal and scoped appropriately

## Future Enhancements

Potential improvements:
- Add filtering options (e.g., only receive specific field updates)
- Implement connection pooling for better scalability
- Add metrics/monitoring for active subscriptions
- Support batch updates to reduce message frequency
