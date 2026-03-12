POST /auth/login
POST /auth/logout

GET /users
POST /users
PATCH /users/:id
DELETE /users/:id

GET /centers
POST /centers
PATCH /centers/:id
DELETE /centers/:id
GET /centers/:centerId/rooms
POST /centers/:centerId/rooms

PATCH /rooms/:id
DELETE /rooms/:id

GET /cpus
POST /cpus
PATCH /cpus/:id
DELETE /cpus/:id

GET /os
POST /os
PATCH /os/:id
DELETE /os/:id

GET /equipment-users
POST /equipment-users
PATCH /equipment-users/:id
DELETE /equipment-users/:id

GET /computers
GET /computers/:id
POST /computers
PATCH /computers/:id
DELETE /computers/:id
GET /computers/:id/audit