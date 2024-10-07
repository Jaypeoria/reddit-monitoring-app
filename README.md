Here’s a comprehensive **documentation** for the installation and usage of the Reddit Monitoring microservices application, which includes OAuth authentication, post fetching, and statistics generation services. The application uses Docker to containerize the services and MongoDB for persistence.

---

## Reddit Monitoring Microservices Application

This application consists of the following microservices:
- **OAuth Service**: Handles Reddit OAuth authentication to acquire access tokens.
- **Fetcher Service**: Fetches Reddit posts from a specific subreddit and stores them in MongoDB.
- **Statistics Service**: Provides statistics on the most upvoted post and user activity based on the posts stored in MongoDB.
- **API Gateway**: Exposes an API to access the posts and statistics.
- **MongoDB**: Stores the Reddit posts for persistence and future analysis.

### Features
- Fetches Reddit posts in near real-time using OAuth authentication.
- Stores posts in MongoDB for persistence.
- Provides statistics such as the most upvoted post and the users with the most posts.
- Microservices architecture with each service running independently.
- Containerized using Docker for easy deployment and scalability.
- Follows SOLID principles.

---

### Prerequisites

- **Docker**: You need to have Docker and Docker Compose installed on your machine. You can download them from [Docker's official website](https://www.docker.com/get-started).
- **Reddit API Credentials**: You need to create a Reddit application to get the necessary OAuth credentials.
  - Go to the [Reddit Apps page](https://www.reddit.com/prefs/apps) and create a new application.
  - You will receive:
    - **Client ID**: A 14-character string displayed under "personal use script".
    - **Client Secret**: A secret key associated with your Reddit app.

---

### Installation

1. **Clone the Repository**

   First, clone the repository to your local machine.

   ```bash
   git clone https://github.com/Jaypeoria/reddit-monitoring-app.git
   cd reddit-monitoring-app
   ```

2. **Set up Environment Variables**

   The application requires several environment variables to be set for Reddit OAuth and MongoDB. These are specified in the `.env` files for each service.

   Create a `.env` file in the root directory of each service (**oauth**, **fetcher**, **statistics**, **gateway**) and add the following variables:

   **Example `.env` file for OAuth Service:**
   ```env
   REDDIT_CLIENT_ID=your-client-id
   REDDIT_CLIENT_SECRET=your-client-secret
   PORT=8081
   ```

   **Example `.env` file for Fetcher and Statistics Services:**
   ```env
   MONGO_URI=mongodb://mongo:27017
   PORT=8082
   ```

   Make sure to replace `your-client-id` and `your-client-secret` with the actual values from your Reddit Developer console.

---

### Running the Application

1. **Build and Start the Application Using Docker Compose**

   Navigate to the root of the project where the `docker-compose.yml` file is located, then run the following command to build and start the application:

   ```bash
   docker-compose up --build
   ```

   Docker Compose will:
   - Build each service (OAuth, Fetcher, Statistics, Gateway).
   - Start all services in the correct order.
   - Spin up a MongoDB instance to store Reddit posts.

   **Note**: The first time you run the build, it might take a few minutes for Docker to pull all necessary images and build the containers.

2. **Verify the Services**

   Once all services are up and running, you can verify the status of the services by accessing the following endpoints:

   - **OAuth Token**: `http://localhost:8081/auth` – This endpoint fetches the OAuth token from Reddit.
   - **Fetcher Service**: `http://localhost:8082/fetch` – This fetches the posts from Reddit and stores them in MongoDB.
   - **Statistics Service**: `http://localhost:8083/stats` – This returns statistics about the most upvoted post and user activity.
   - **API Gateway**: `http://localhost:8080/posts` and `http://localhost:8080/stats` – These endpoints provide access to posts and statistics.

---

### Usage

Once the application is running, you can interact with the API using tools like **cURL** or **Postman**.

1. **Fetch OAuth Token**

   The OAuth Service provides a token to interact with the Reddit API. Make a GET request to the `/auth` endpoint:

   ```bash
   curl http://localhost:8081/auth
   ```

   Example Response:
   ```json
   {
     "access_token": "your-reddit-oauth-token"
   }
   ```

2. **Fetch Reddit Posts**

   The Fetcher Service fetches Reddit posts from a specific subreddit (by default, it’s set to `golang`). You can trigger fetching by making a GET request to `/fetch`:

   ```bash
   curl http://localhost:8082/fetch
   ```

   Example Response:
   ```json
   {
     "posts": [
       {
         "author": "author_name",
         "title": "post_title",
         "score": 123
       },
       ...
     ]
   }
   ```

3. **Get Statistics**

   The Statistics Service provides insights into the posts stored in MongoDB. You can retrieve the statistics using the `/stats` endpoint:

   ```bash
   curl http://localhost:8083/stats
   ```

   Example Response:
   ```json
   {
     "most_upvoted_post": {
       "author": "author_name",
       "title": "most_upvoted_post_title",
       "score": 999
     },
     "user_post_counts": {
       "author_name1": 10,
       "author_name2": 8
     }
   }
   ```

4. **Use API Gateway**

   The API Gateway aggregates the data from both the **Fetcher** and **Statistics** services. Use the following endpoints to interact with the gateway:

   - **Fetch Posts**: `http://localhost:8080/posts`
   - **Get Statistics**: `http://localhost:8080/stats`

---

### Stopping the Application

To stop the application, press `CTRL + C` in the terminal where Docker Compose is running. Alternatively, you can stop and remove the containers using:

```bash
docker-compose down
```

This will stop all running containers and clean up the environment.

---

### Monitoring and Logs

1. **View Logs**

   To view the logs for each service, you can use the following Docker command:

   ```bash
   docker-compose logs -f
   ```

   This will show the logs for all running services. You can also specify a specific service (e.g., fetcher-service) to see only that service's logs:

   ```bash
   docker-compose logs -f fetcher-service
   ```

