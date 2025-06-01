from pydantic import BaseModel
from uuid import UUID
from datetime import datetime
from typing import Optional

class Follower(BaseModel):
    user_id: UUID
    follower_id: UUID
    followed_at: Optional[datetime] = None

class Followee(BaseModel):
    user_id: UUID
    followee_id: UUID
    followed_at: Optional[datetime] = None
