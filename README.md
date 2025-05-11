


1. `main.go`: This is the entry point of the application, where the main function is defined. The main function creates a new instance of the `App` struct, which contains references to the various components of the application.
2. `app.go`: This is the central struct that holds references to all the other components of the application. It provides a way to initialize and manage the application's state.
3. `ratelimiter.go`: This file defines a rate limiter component, which is used to control the number of requests that can be made to the API within a certain time frame. The `Limiter` type defines the interface for rate limiting, and the `FixedWindowLimiter` struct implements this interface with a fixed window size.
4. `store.go`: This file defines various types for storing data in a database, including `User`, `Post`, `Comment`, and `Follow`. It also includes functions for creating, reading, updating, and deleting data in the database.
5. `users.go`: This file defines the `User` type and provides functions for creating, reading, updating, and deleting users in the database.
6. `posts.go`: This file defines the `Post` type and provides functions for creating, reading, updating, and deleting posts in the database.
7. `comments.go`: This file defines the `Comment` type and provides functions for creating, reading, updating, and deleting comments on posts in the database.
8. `follow.go`: This file defines the `Follow` type and provides functions for following and unfollowing users in the database.
9. `main.Controller`: This is a controller function that handles incoming HTTP requests and calls the appropriate functions to handle the request.
10. `models`: This folder contains the structs and interfaces for the application's data models, including `User`, `Post`, `Comment`, and `Follow`.
11. `views`: This folder contains the HTML templates for the application's user interface.
12. `routes`: This folder contains the routes for the application's HTTP requests, which are defined using the `main.Controller` function.

