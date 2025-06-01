from fastapi import APIRouter, HTTPException
from uuid import UUID, uuid4
from app.core.cassandra_connector import CassandraConnector
from app.models.models_post import PostByUser, PostByDateBucket, PostByKeyword
from typing import List
from datetime import datetime

router = APIRouter()
cassandra = CassandraConnector()
session = cassandra.connect()

@router.get("/posts/user/{user_id}", response_model=List[PostByUser])
def get_posts_by_user(user_id: UUID, limit: int = 10):
    query = "SELECT * FROM posts_by_user WHERE user_id = %s LIMIT %s"
    results = session.execute(query, (user_id, limit))
    return [PostByUser(**row._asdict()) for row in results]

@router.get("/posts/date/{date_bucket}", response_model=List[PostByDateBucket])
def get_posts_by_date_bucket(date_bucket: str, limit: int = 10):
    query = "SELECT * FROM posts_by_date_bucket WHERE date_bucket = %s LIMIT %s"
    results = session.execute(query, (date_bucket, limit))
    return [PostByDateBucket(**row._asdict()) for row in results]

@router.get("/posts/keyword/{keyword}", response_model=List[PostByKeyword])
def get_posts_by_keyword(keyword: str, limit: int = 10):
    query = "SELECT * FROM posts_by_keyword WHERE keyword = %s LIMIT %s"
    results = session.execute(query, (keyword, limit))
    return [PostByKeyword(**row._asdict()) for row in results]

@router.post("/posts", response_model=PostByUser)
def create_post(post: PostByUser):
    query = """
    INSERT INTO posts_by_user (user_id, created_at, post_id, title, description, image_urls)
    VALUES (%s, %s, %s, %s, %s, %s)
    """
    post_id = post.post_id or uuid4()
    created_at = post.created_at or datetime.utcnow()
    session.execute(query, (
        post.user_id,
        created_at,
        post_id,
        post.title,
        post.description,
        post.image_urls
    ))
    return PostByUser(
        user_id=post.user_id,
        created_at=created_at,
        post_id=post_id,
        title=post.title,
        description=post.description,
        image_urls=post.image_urls
    )

@router.delete("/posts/user/{user_id}/{post_id}", response_model=PostByUser)
def delete_post_by_user(user_id: UUID, post_id: UUID):
    # Busca el post
    result = session.execute(
        "SELECT * FROM posts_by_user WHERE user_id = %s AND post_id = %s ALLOW FILTERING",
        (user_id, post_id)
    ).one()
    if not result:
        raise HTTPException(status_code=404, detail="Post not found")
    # Elimina el post
    session.execute(
        "DELETE FROM posts_by_user WHERE user_id = %s AND created_at = %s AND post_id = %s",
        (user_id, result.created_at, post_id)
    )
    return PostByUser(**result._asdict())

@router.delete("/posts/date/{date_bucket}/{post_id}", response_model=PostByDateBucket)
def delete_post_by_date_bucket(date_bucket: str, post_id: UUID):
    result = session.execute(
        "SELECT * FROM posts_by_date_bucket WHERE date_bucket = %s AND post_id = %s ALLOW FILTERING",
        (date_bucket, post_id)
    ).one()
    if not result:
        raise HTTPException(status_code=404, detail="Post not found")
    session.execute(
        "DELETE FROM posts_by_date_bucket WHERE date_bucket = %s AND created_at = %s AND post_id = %s",
        (date_bucket, result.created_at, post_id)
    )
    return PostByDateBucket(**result._asdict())

@router.delete("/posts/keyword/{keyword}/{post_id}", response_model=PostByKeyword)
def delete_post_by_keyword(keyword: str, post_id: UUID):
    result = session.execute(
        "SELECT * FROM posts_by_keyword WHERE keyword = %s AND post_id = %s ALLOW FILTERING",
        (keyword, post_id)
    ).one()
    if not result:
        raise HTTPException(status_code=404, detail="Post not found")
    session.execute(
        "DELETE FROM posts_by_keyword WHERE keyword = %s AND created_at = %s AND post_id = %s",
        (keyword, result.created_at, post_id)
    )
    return PostByKeyword(**result._asdict())
