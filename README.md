# Go Blog
#### Description:
##### Features
- Create Blog Posts: Users can create blog posts with a title, content, and an expiration date.
- View Blog Posts: View individual blog posts with detailed information, including creation and expiration dates.
- Automatic Expiration: Blog posts can be set to automatically expire at a specified date.
- Connection Pooling: Efficient database connection pooling using Go’s database/sql package.
- Template Rendering: Dynamic HTML rendering with Go templates for serving content to the frontend.
##### Technologies Used
- Go: Language used for building the backend application.
- MySQL: Database for storing blog posts.
- HTML Templates: For dynamically rendering blog post data to the web pages.
- Postman: Recommended tool for testing API endpoints.
##### Project Structure
- main.go: Entry point of the application, handles server initialization and routing.
- models: Contains the BlogModel struct, which defines database interactions for blog posts.
- templates: Holds HTML template files for rendering views.
- handlers.go: Defines HTTP handler functions for different routes.
- templates.go: Template management and caching logic.