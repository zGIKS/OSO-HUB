from fastapi import FastAPI
from app.routers import users, posts, comments, likes, follows, feed

app = FastAPI()

app.include_router(users.router)
app.include_router(posts.router)
app.include_router(comments.router)
app.include_router(likes.router)
app.include_router(follows.router)
app.include_router(feed.router)

if __name__ == "__main__":
    pass
