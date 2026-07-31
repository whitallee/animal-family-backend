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

## Infra / DevOps

- [ ] Verify `OPENAI_API_KEY` has been added to the live ECS task definition (`brindl-backend`) — added to local `.env` and tested locally on 2026-07-27 but not yet confirmed deployed
- [ ] When setting up Terraform/IaC for ECS (see CLAUDE.md "Deployment"): move secrets out of the task definition's plaintext `environment` block into AWS Secrets Manager/SSM, referenced via the task def's `secrets` field. Terraform should declare the full set of secret references as code so a re-render can't silently drop one, and so secrets aren't readable in plaintext via `describe-task-definition`

## Features

- [ ] Action History feature
- [ ] Consider multiple subjects per task (e.g., feed all 4 ferrets as one task instead of per-enclosure)
- [ ] Permanent pet ownership transfer (request/accept flow between users)
- [ ] Temporary ownership transfer for pet sitters (time-bound access with configurable permissions)
