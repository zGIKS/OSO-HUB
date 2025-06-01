# models.py
# backend/app/models.py
from pydantic import BaseModel, EmailStr
from typing import Optional
from datetime import datetime
from uuid import UUID

class UserCreate(BaseModel):
    username: str
    email: EmailStr
    password_hash: str
    full_name: Optional[str] = None
    bio: Optional[str] = None
    avatar_url: Optional[str] = None

class UserDB(BaseModel):
    user_id: UUID
    username: str
    email: EmailStr
    created_at: datetime

class UserUpdate(BaseModel):
    username: Optional[str] = None
    email: Optional[EmailStr] = None

class User(UserCreate):
    user_id: UUID
    created_at: datetime
