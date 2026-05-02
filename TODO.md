# To-Do

## In Progress

- [ ] Task endpoints (routes and store)

## Backend

- [ ] Add `CreateAnimalAndEnclosure` for simultaneous creation (`CreateEnclosureWithAnimals` already exists)
- [ ] Add `UpdateUser` function and route
- [ ] Add `UpdateAnimalSubject` and `UpdateEnclosureSubject` functions and routes
- [ ] Add ownership transfer request flow (user must accept before ownership changes)
  - `handleUserUpdateAnimalOwner`
  - `handleUserUpdateEnclosureOwner`
  - `handleUserUpdateTaskOwner`
- [ ] Add duplicate check when changing ownership
- [ ] Fix transaction rollbacks in Task service (`CreateTask`, `DeleteTaskById`)
- [ ] Modularize repeated ownership checks across route handlers
- [ ] Use goroutines/WaitGroups for concurrent batch operations (e.g., `DeleteUserById` loops)
  - Will require more modular store functions

## Features

- [ ] Action History feature
- [ ] Consider multiple subjects per task (e.g., feed all 4 ferrets as one task instead of per-enclosure)
- [ ] Permanent pet ownership transfer (request/accept flow between users)
- [ ] Temporary ownership transfer for pet sitters (time-bound access with configurable permissions)
