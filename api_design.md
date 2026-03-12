# API Design

## Authorization

Requests are authenticated via a Bearer token obtained from `POST /auth/login`.

Role hierarchy (least → most privileged):
- **readonly** – read-only access
- **tutor** – manages their own class, students, and assignments
- **editor** – manages inventory data (create/update most resources)
- **admin** – full access including user and role management

---

## Auth

```
POST /auth/login     – public, returns token
POST /auth/logout    – any authenticated user
```

---

## Users

```
GET    /users        – any authenticated
POST   /users        – admin
PATCH  /users/{id}   – admin
DELETE /users/{id}   – admin
```

---

## Roles

```
GET    /roles        – any authenticated
POST   /roles        – admin
DELETE /roles/{id}   – admin
```

---

## Centers

```
GET    /centers        – any authenticated
POST   /centers        – editor, admin
PATCH  /centers/{id}   – editor, admin
DELETE /centers/{id}   – admin
```

---

## Rooms

```
GET    /centers/{centerId}/rooms   – any authenticated
POST   /centers/{centerId}/rooms   – editor, admin
PATCH  /rooms/{id}                 – editor, admin
DELETE /rooms/{id}                 – admin
```

---

## CPUs

```
GET    /cpus        – any authenticated
POST   /cpus        – editor, admin
PATCH  /cpus/{id}   – editor, admin
DELETE /cpus/{id}   – admin
```

---

## Operating Systems

```
GET    /os        – any authenticated
POST   /os        – editor, admin
PATCH  /os/{id}   – editor, admin
DELETE /os/{id}   – admin
```

---

## Equipment Users

```
GET    /equipment-users        – any authenticated
POST   /equipment-users        – editor, admin
PATCH  /equipment-users/{id}   – editor, admin
DELETE /equipment-users/{id}   – admin
```

---

## Brands

```
GET    /brands        – any authenticated
POST   /brands        – editor, admin
PATCH  /brands/{id}   – editor, admin
DELETE /brands/{id}   – admin
```

---

## Laptop Models

```
GET    /laptop-models           – any authenticated
POST   /laptop-models           – editor, admin
GET    /laptop-models/{id}      – any authenticated
PATCH  /laptop-models/{id}      – editor, admin
DELETE /laptop-models/{id}      – admin
```

---

## Desktop Models

```
GET    /desktop-models          – any authenticated
POST   /desktop-models          – editor, admin
GET    /desktop-models/{id}     – any authenticated
PATCH  /desktop-models/{id}     – editor, admin
DELETE /desktop-models/{id}     – admin
```

---

## Computers (base)

Generic computer listing and deletion. Use `/desktops` or `/laptops` for typed access.

```
GET    /computers        – any authenticated; returns all computers with type field
GET    /computers/{id}   – any authenticated; returns base fields
DELETE /computers/{id}   – editor, admin; cascades to desktop/laptop subtype
```

---

## Desktops

```
GET    /desktops        – any authenticated; full join (computer + desktop)
POST   /desktops        – editor, admin
GET    /desktops/{id}   – any authenticated
PATCH  /desktops/{id}   – editor, admin
```

### POST / PATCH /desktops body fields

| Field             | Type    | Required (POST) | Notes                               |
|-------------------|---------|-----------------|-------------------------------------|
| hostname          | string  | yes             |                                     |
| room_id           | integer | no              |                                     |
| observations      | string  | no              |                                     |
| desktop_model_id  | integer | no              | null → fully manual spec            |
| cpu_id            | integer | no              |                                     |
| ram_gb            | integer | no              |                                     |
| ram_type          | string  | no              | DDR3 / DDR4 / DDR5 / None           |
| storage_gb        | integer | no              |                                     |
| storage_type      | string  | no              | HDD / SSD / NVMe / None             |
| os_id             | integer | no              |                                     |
| equipment_user_id | integer | no              |                                     |
| has_wifi_card     | boolean | yes (POST)      |                                     |
| mac_address       | string  | no              | required if has_wifi_card = true    |

---

## Laptops

```
GET    /laptops        – any authenticated; full join (computer + laptop + model)
POST   /laptops        – editor, admin
GET    /laptops/{id}   – any authenticated
PATCH  /laptops/{id}   – editor, admin
```

### POST / PATCH /laptops body fields

| Field             | Type    | Required (POST) | Notes                                   |
|-------------------|---------|-----------------|------------------------------------------|
| hostname          | string  | yes             |                                          |
| room_id           | integer | no              |                                          |
| observations      | string  | no              |                                          |
| laptop_model_id   | integer | yes             |                                          |
| ram_gb            | integer | no              | null → inherit from model                |
| ram_type          | string  | no              | DDR3 / DDR4 / DDR5 / None; null = model |
| storage_gb        | integer | no              | null → inherit from model                |
| storage_type      | string  | no              | HDD / SSD / NVMe / None; null = model   |
| mac_address       | string  | no              |                                          |
| os_id             | integer | no              | null → inherit from model                |
| equipment_user_id | integer | no              | null → student laptop                    |

---

## Laptop Assignments

```
GET    /laptops/{laptopId}/assignments   – any authenticated
POST   /laptops/{laptopId}/assignments   – editor, admin, tutor (own class only)
GET    /assignments/{id}                 – any authenticated
PATCH  /assignments/{id}                 – editor, admin, tutor (own class only)
DELETE /assignments/{id}                 – editor, admin, tutor (own class only)
```

### POST / PATCH /assignments body fields

| Field         | Type   | Required (POST) | Notes                              |
|---------------|--------|-----------------|------------------------------------|
| student_id    | integer | yes            |                                    |
| class_id      | integer | yes            | snapshot of student's class        |
| academic_year | string  | yes            | format: `YYYY-YYYY` e.g. `2024-2025` |

---

## Cycles

```
GET    /cycles        – any authenticated
POST   /cycles        – editor, admin
PATCH  /cycles/{id}   – editor, admin
DELETE /cycles/{id}   – admin
```

---

## Classes

```
GET    /cycles/{cycleId}/classes   – any authenticated
POST   /cycles/{cycleId}/classes   – editor, admin
GET    /classes/{id}               – any authenticated
PATCH  /classes/{id}               – editor, admin, tutor (own class only)
DELETE /classes/{id}               – admin
```

### POST /cycles/{cycleId}/classes body fields

| Field             | Type    | Required | Notes                               |
|-------------------|---------|----------|-------------------------------------|
| course            | integer | yes      | e.g. 1, 2                           |
| class_label       | string  | no       | default `A`                         |
| shift             | string  | yes      | `morning` / `afternoon`             |
| tutor_app_user_id | integer | no       |                                     |

---

## Students

```
GET    /classes/{classId}/students   – any authenticated
POST   /classes/{classId}/students   – editor, admin, tutor (own class only)
GET    /students/{id}                – any authenticated
PATCH  /students/{id}                – editor, admin, tutor (own class only)
DELETE /students/{id}                – admin
```

---

## Audit Log

```
GET /audit?table=<tableName>&record_id=<id>   – admin
```

Returns the audit history for a specific record in chronological descending order. Supported table names: `computer`, `desktop`, `laptop`.
