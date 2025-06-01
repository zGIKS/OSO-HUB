from pydantic import BaseModel
from uuid import UUID
from datetime import datetime
from typing import Optional

class Comment(BaseModel):
    post_id: UUID
    created_at: datetime
    comment_id: UUID
    commenter_id: UUID
    content: str
