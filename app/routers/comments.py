from fastapi import APIRouter, HTTPException
from uuid import UUID, uuid1
from app.core.cassandra_connector import CassandraConnector
from app.models.models_comment import Comment
from typing import List
from datetime import datetime

router = APIRouter()
cassandra = CassandraConnector()
session = cassandra.connect()

@router.get("/comments/{post_id}", response_model=List[Comment])
def get_comments_by_post(post_id: UUID, limit: int = 10):
    query = "SELECT * FROM comments_by_post WHERE post_id = %s LIMIT %s"
    results = session.execute(query, (post_id, limit))
    return [Comment(**row._asdict()) for row in results]

@router.post("/comments", response_model=Comment)
def create_comment(comment: Comment):
    query = """
    INSERT INTO comments_by_post (post_id, created_at, comment_id, commenter_id, content)
    VALUES (%s, %s, %s, %s, %s)
    """
    comment_id = comment.comment_id or uuid1()
    created_at = comment.created_at or datetime.utcnow()
    session.execute(query, (
        comment.post_id,
        created_at,
        comment_id,
        comment.commenter_id,
        comment.content
    ))
    return Comment(
        post_id=comment.post_id,
        created_at=created_at,
        comment_id=comment_id,
        commenter_id=comment.commenter_id,
        content=comment.content
    )

@router.put("/comments/{post_id}/{comment_id}", response_model=Comment)
def update_comment(post_id: UUID, comment_id: UUID, content: str):
    # Busca el comentario
    result = session.execute(
        "SELECT * FROM comments_by_post WHERE post_id = %s AND comment_id = %s",
        (post_id, comment_id)
    ).one()
    if not result:
        raise HTTPException(status_code=404, detail="Comment not found")
    # Actualiza el contenido
    session.execute(
        "UPDATE comments_by_post SET content = %s WHERE post_id = %s AND created_at = %s AND comment_id = %s",
        (content, post_id, result.created_at, comment_id)
    )
    # Devuelve el comentario actualizado
    updated = result._asdict()
    updated["content"] = content
    return Comment(**updated)
