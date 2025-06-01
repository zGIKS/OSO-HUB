from pydantic import BaseModel
from uuid import UUID
from datetime import datetime
from typing import List, Optional

class PostByUser(BaseModel):
    user_id: UUID
    created_at: datetime
    post_id: UUID
    title: str
    description: Optional[str] = None
    image_urls: Optional[List[str]] = None

class PostByDateBucket(BaseModel):
    date_bucket: str
    created_at: datetime
    post_id: UUID
    user_id: UUID
    title: str
    description: Optional[str] = None
    image_urls: Optional[List[str]] = None

class PostByKeyword(BaseModel):
    keyword: str
    created_at: datetime
    post_id: UUID
    user_id: UUID
