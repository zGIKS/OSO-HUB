from fastapi import APIRouter, HTTPException
from uuid import UUID
from app.core.cassandra_connector import CassandraConnector
from app.models.models_follow import Follower, Followee
from typing import List
from datetime import datetime

router = APIRouter()
cassandra = CassandraConnector()
session = cassandra.connect()

@router.get("/followers/{user_id}", response_model=List[Follower])
def get_followers(user_id: UUID, limit: int = 10):
    query = "SELECT * FROM followers_by_user WHERE user_id = %s LIMIT %s"
    results = session.execute(query, (user_id, limit))
    return [Follower(**row._asdict()) for row in results]

@router.get("/followees/{user_id}", response_model=List[Followee])
def get_followees(user_id: UUID, limit: int = 10):
    query = "SELECT * FROM followees_by_user WHERE user_id = %s LIMIT %s"
    results = session.execute(query, (user_id, limit))
    return [Followee(**row._asdict()) for row in results]

@router.post("/followers", response_model=Follower)
def create_follower(follower: Follower):
    query = """
    INSERT INTO followers_by_user (user_id, follower_id, followed_at)
    VALUES (%s, %s, %s)
    """
    followed_at = follower.followed_at or datetime.utcnow()
    session.execute(query, (
        follower.user_id,
        follower.follower_id,
        followed_at
    ))
    return Follower(
        user_id=follower.user_id,
        follower_id=follower.follower_id,
        followed_at=followed_at
    )

@router.post("/followees", response_model=Followee)
def create_followee(followee: Followee):
    query = """
    INSERT INTO followees_by_user (user_id, followee_id, followed_at)
    VALUES (%s, %s, %s)
    """
    followed_at = followee.followed_at or datetime.utcnow()
    session.execute(query, (
        followee.user_id,
        followee.followee_id,
        followed_at
    ))
    return Followee(
        user_id=followee.user_id,
        followee_id=followee.followee_id,
        followed_at=followed_at
    )

@router.delete("/followers", response_model=Follower)
def delete_follower(user_id: UUID, follower_id: UUID):
    result = session.execute(
        "SELECT * FROM followers_by_user WHERE user_id = %s AND follower_id = %s",
        (user_id, follower_id)
    ).one()
    if not result:
        raise HTTPException(status_code=404, detail="Follower not found")
    session.execute(
        "DELETE FROM followers_by_user WHERE user_id = %s AND follower_id = %s",
        (user_id, follower_id)
    )
    return Follower(**result._asdict())

@router.delete("/followees", response_model=Followee)
def delete_followee(user_id: UUID, followee_id: UUID):
    result = session.execute(
        "SELECT * FROM followees_by_user WHERE user_id = %s AND followee_id = %s",
        (user_id, followee_id)
    ).one()
    if not result:
        raise HTTPException(status_code=404, detail="Followee not found")
    session.execute(
        "DELETE FROM followees_by_user WHERE user_id = %s AND followee_id = %s",
        (user_id, followee_id)
    )
    return Followee(**result._asdict())
