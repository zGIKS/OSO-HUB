from pydantic import BaseModel
from uuid import UUID
from datetime import datetime
from typing import Optional

class LikeCount(BaseModel):
    post_id: UUID
    likes: int

class LikeByPost(BaseModel):
    post_id: UUID
    user_id: UUID
    liked_at: Optional[datetime] = None
