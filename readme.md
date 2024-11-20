# **Forum Web Application**

## **Project Description**
This project consists of building a **web-based forum** that facilitates:
- **User communication**.
- **Post categorization**.
- **Like/dislike functionality**.
- **Post filtering**.

It leverages **SQLite** for database management, implements **user authentication**, and supports a **Dockerized environment** for streamlined deployment.

---

## **Features**

### **1. User Authentication**

#### **User Registration**
- Users must provide:
  - **Email** (must be unique).
  - **Username**.
  - **Password**.
- Password encryption is implemented using **bcrypt** (*Bonus*).

#### **Login System**
- Users log in with their **email and password**.
- Sessions are managed using **cookies**, which include:
  - A **unique session ID**.
  - An **expiration date**.
- Users can have only **one active session** at a time.

---

### **2. Posts and Comments**

#### **Creating Posts**
- Registered users can **create posts**.
- Posts can be tagged with **one or more categories** for better organization.

#### **Creating Comments**
- Registered users can **comment on posts**.
- Comments are displayed **below their respective posts**.

#### **Viewing Content**
- All users (**registered or not**) can **view posts and comments**.

---

### **3. Likes and Dislikes**
- **Registered users** can:
  - **Like** or **dislike** posts and comments.
- The **total number of likes and dislikes** is visible to everyone.

---

### **4. Filters**
Users can filter posts by:
1. **Categories**: View posts within specific subforums.
2. **Created Posts**: View posts made by the logged-in user (**for registered users only**).
3. **Liked Posts**: View posts liked by the logged-in user (**for registered users only**).

---

### **5. Error Handling**
The application handles errors gracefully with:
- **User Input Validation**:
  - Example: Duplicate email returns a **400 Bad Request**.
- **Database Errors**:
  - Issues like connection failures are logged and handled appropriately.
- **Unexpected Errors**:
  - A **500 Internal Server Error** is returned with a descriptive message.

---

## **SQLite Integration**
The database includes the following tables:
- **Users**: Stores email, username, and hashed passwords.
- **Posts**: Includes title, content, categories, and author.
- **Comments**: Tracks content, associated post, and author.
- **Likes/Dislikes**: Logs user interactions with posts and comments.

### **Example SQL Queries**
- **CREATE**: Create the necessary tables.
- **INSERT**: Add users, posts, comments, likes/dislikes.
- **SELECT**: Retrieve posts, comments, and user interactions.

---

## **Docker Setup**

### **1. Building the Docker Image**
To create a Docker image for this application, run:
```bash
docker build -t forum-web-app .
