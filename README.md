# Go Social Media API

A simple yet robust social media backend written in **Go** using the **Chi router**, **PostgreSQL**, **JWT authentication**, and **role-based authorization**. The project follows the repository pattern, includes rate limiting, CI/CD via GitHub Actions, Swagger documentation, Redis caching, and graceful shutdowns.

---

## 🚀 Features

- User registration and authentication (JWT)
- Role-based authorization (`user`, `moderator`, `admin`)
- Create, update, delete posts (with ownership checks)
- Post comments
- Follow/unfollow users
- User feeds
- Redis caching
- Graceful shutdowns
- Swagger API docs
- Health and debug endpoints
- Secure routes with middleware
- CI/CD using GitHub Actions

---

## 📦 Base URL

All routes are versioned under:  
```

/api/v1

```

---

## 📚 API Endpoints

### 🔐 Authentication

| Method | Endpoint                 | Description         |
|--------|--------------------------|---------------------|
| POST   | `/authenticate/register` | Register a new user |
| POST   | `/authenticate/login`    | Login a user        |

---

### 📝 Posts

Requires JWT via `AuthTokenMiddleware`.

| Method | Endpoint                          | Description                     |
|--------|-----------------------------------|---------------------------------|
| POST   | `/post/`                          | Create a new post               |
| POST   | `/post/comment`                   | Add a comment to a post         |
| GET    | `/post/{postID}`                  | Retrieve a post by ID           |
| PATCH  | `/post/{postID}`                  | Update post (requires ownership or `moderator` role) |
| DELETE | `/post/{postID}`                  | Delete post (requires `admin`)  |

---

### 👤 User

| Method | Endpoint                        | Description                        |
|--------|---------------------------------|------------------------------------|
| PUT    | `/user/activate/{token}`        | Activate user account              |
| GET    | `/user/{userID}`                | Get user profile (auth required)   |
| PUT    | `/user/{userID}/follow`         | Follow a user (auth required)      |
| PUT    | `/user/{userID}/unfollow`       | Unfollow a user (auth required)    |

---

### 📰 Feed

| Method | Endpoint          | Description                  |
|--------|-------------------|------------------------------|
| GET    | `/user/feed`      | Get user’s social media feed |

---

### 🛠️ Debug & Health

Protected by basic auth.

| Method | Endpoint        | Description              |
|--------|-----------------|--------------------------|
| GET    | `/debug/vars`   | App debug variables (expvar) |
| GET    | `/health`       | Health check endpoint    |

---

### 📄 Swagger Documentation

Interactive API docs are available at:

```

/api/v1/swagger/index.html

````

---

## 🧰 Tech Stack

- **Go** (Golang)
- **Chi** router
- **PostgreSQL**
- **Redis** (cache)
- **JWT** auth
- **Zap** (logging)
- **Swagger** (docs)
- **GitHub Actions** (CI/CD)

---

## 🛡️ Security

- JWT-based authentication
- Role-based access control
- Middleware chain with auth checks
- Graceful error handling
- Secure routes (basic auth for internal endpoints)

---


## ✅ TODO (Future Enhancements)

* Notifications (e.g., for comments or follows)
* gRPC-based internal microservices
* Media/image uploads
* WebSockets for real-time updates
* Rate limiter backed by Redis
* Admin dashboard or analytics endpoints

---

## 📬 Feedback or Issues?

Feel free to open an issue or PR!





