# Animal Family Backend
This is very much a work in progress right now. Let's call it a Pre-Alpha-Alpha. If you'd like to collab in any way, please reach out! Find my contact info on [whitcodes.dev/contact](whitcodes.dev/contact). If you'd like to check out my first prototype, I've got it hosted on Vercel currently at [this](animal-family.vercel.app) link. Also the text-notification functionality is down, but I'll be working on that once I have this new backend up and running.

## To-Do List
- **WIP** Add Get functions
- Add checks for creation endpoints to see if subject already exists
    - handleCreateHabitat: check by name
    - handleCreateSpecies: check by comName
    - handleCreateSpecies: check by sciName
- Add Update functions
- Implement the Tasks Feature
- Use GoRoutines for concurrent async requests like for DeleteUserById animal and enclosure requests or any for loop db requests
- Implement Action History Feature
- eventually add more admin routes such as DeleteAnimalById or DeleteEnclosureById with an admin handler

## Entity Relationship Diagram
[Here's a diagram](https://docs.google.com/drawings/d/1Vi1yngr4CeXXt-slRGJsLI35_R-y-oIHlZ466be_wx8/edit?usp=sharing) that I made of the DB schema. Feel free to leave comments on the Drawing.

## Plans for Frontend Applications
- Use a cache heavily and store all user animals, enclosures, and tasks
    - invalidate and refetch cache every 15 minutes
    - invalidate and refetch cache if user updates any of their animals, enclosures, or tasks
    - should improve performance massively compared to my v1 Next.js frontend