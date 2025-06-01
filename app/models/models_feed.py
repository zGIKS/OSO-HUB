from pydantic import BaseModel
from uuid import UUID
from datetime import datetime

class FeedByUser(BaseModel):
    user_id: UUID
    post_created_at: datetime
    post_id: UUID
    author_id: UUID
    title: str
