from fastapi import APIRouter, HTTPException
from uuid import UUID
from app.core.cassandra_connector import CassandraConnector
from app.models.models_like import LikeCount, LikeByPost
from typing import List
from datetime import datetime

router = APIRouter()
cassandra = CassandraConnector()
session = cassandra.connect()

@router.get("/likes/count/{post_id}", response_model=LikeCount)
def get_likes_count(post_id: UUID):
    query = "SELECT * FROM likes_count WHERE post_id = %s"
    result = session.execute(query, (post_id,)).one()
    if not result:
        raise HTTPException(status_code=404, detail="Post not found")
    return LikeCount(**result._asdict())

@router.get("/likes/{post_id}", response_model=List[LikeByPost])
def get_likes_by_post(post_id: UUID, limit: int = 10):
    query = "SELECT * FROM likes_by_post WHERE post_id = %s LIMIT %s"
    results = session.execute(query, (post_id, limit))
    return [LikeByPost(**row._asdict()) for row in results]

@router.post("/likes", response_model=LikeByPost)
def create_like(like: LikeByPost):
    query = """
    INSERT INTO likes_by_post (post_id, user_id, liked_at)
    VALUES (%s, %s, %s)
    """
    liked_at = like.liked_at or datetime.utcnow()
    session.execute(query, (
        like.post_id,
        like.user_id,
        liked_at
    ))
    # Incrementa el contador de likes
    session.execute("UPDATE likes_count SET likes = likes + 1 WHERE post_id = %s", (like.post_id,))
    return LikeByPost(
        post_id=like.post_id,
        user_id=like.user_id,
        liked_at=liked_at
    )

@router.delete("/likes", response_model=LikeByPost)
def delete_like(post_id: UUID, user_id: UUID):
    # Verifica si existe el like
    result = session.execute(
        "SELECT * FROM likes_by_post WHERE post_id = %s AND user_id = %s",
        (post_id, user_id)
    ).one()
    if not result:
        raise HTTPException(status_code=404, detail="Like not found")
    # Elimina el like
    session.execute(
        "DELETE FROM likes_by_post WHERE post_id = %s AND user_id = %s",
        (post_id, user_id)
    )
    # Decrementa el contador de likes
    session.execute(
        "UPDATE likes_count SET likes = likes - 1 WHERE post_id = %s",
        (post_id,)
    )
    return LikeByPost(**result._asdict())
