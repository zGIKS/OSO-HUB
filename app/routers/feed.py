from fastapi import APIRouter, HTTPException
from uuid import UUID
from app.core.cassandra_connector import CassandraConnector
from app.models.models_feed import FeedByUser
from typing import List

router = APIRouter()
cassandra = CassandraConnector()
session = cassandra.connect()

@router.get("/feed/{user_id}", response_model=List[FeedByUser])
def get_feed_by_user(user_id: UUID, limit: int = 20):
    query = "SELECT * FROM feed_by_user WHERE user_id = %s LIMIT %s"
    results = session.execute(query, (user_id, limit))
    return [FeedByUser(**row._asdict()) for row in results]

@router.delete("/feed/{user_id}/{post_id}", response_model=FeedByUser)
def delete_feed_item(user_id: UUID, post_id: UUID):
    result = session.execute(
        "SELECT * FROM feed_by_user WHERE user_id = %s AND post_id = %s ALLOW FILTERING",
        (user_id, post_id)
    ).one()
    if not result:
        raise HTTPException(status_code=404, detail="Feed item not found")
    session.execute(
        "DELETE FROM feed_by_user WHERE user_id = %s AND post_created_at = %s AND post_id = %s",
        (user_id, result.post_created_at, post_id)
    )
    return FeedByUser(**result._asdict())
