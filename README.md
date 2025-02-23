# Animal Family Backend
This is very much a work in progress right now. Let's call it a Pre-Alpha-Alpha. If you'd like to collab in any way, please reach out! Find my contact info on [whitcodes.dev/contact](whitcodes.dev/contact). If you'd like to check out my first prototype, I've got it hosted on Vercel currently at [this](animal-family.vercel.app) link. Also the text-notification functionality is down, but I'll be working on that once I have this new backend up and running.

## To-Do List
- Implement delete functions and handlers
    - DeleteSpeciesByIdAsAdmin
    - DeleteHabitatByIdAsAdmin

- Change functions to have ...AsAdmin functions instead
    - CreateHabitat
    - CreateSpecies
    - GetAnimals
    - GetEnclosures

- Update all transaction functions to match the structure of DeleteAnimalByIdWithUserId with error handling

- Add checks for creation endpoints to see if subject already exists

- Add Update functions

- Implement the Tasks Feature

- Implement Action History Feature

## Entity Relationship Diagram
[Here's a diagram](https://docs.google.com/drawings/d/1Vi1yngr4CeXXt-slRGJsLI35_R-y-oIHlZ466be_wx8/edit?usp=sharing) that I made of the DB schema. Feel free to leave comments on the Drawing.