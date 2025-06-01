# users.py

from fastapi import APIRouter, HTTPException, Body
from uuid import UUID, uuid4
from app.core.cassandra_connector import CassandraConnector
from app.models.models import User, UserCreate
from typing import List
from datetime import datetime

router = APIRouter()

# Instancia global del conector (en producción usaría dependencia)
cassandra = CassandraConnector()
session = cassandra.connect()

@router.get("/users/{user_id}", response_model=User)
def get_user(user_id: UUID):
    query = "SELECT * FROM users_by_id WHERE user_id = %s"
    result = session.execute(query, (user_id,)).one()
    if not result:
        raise HTTPException(status_code=404, detail="User not found")
    return User(**result._asdict())

@router.get("/users", response_model=List[User])
def list_users(limit: int = 10):
    query = f"SELECT * FROM users_by_id LIMIT {limit}"
    results = session.execute(query)
    return [User(**row._asdict()) for row in results]

@router.post("/users", response_model=User)
def create_user(user: UserCreate):
    user_id = uuid4()
    created_at = datetime.utcnow()
    query = """
    INSERT INTO users_by_id (user_id, username, email, password_hash, full_name, bio, avatar_url, created_at)
    VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
    """
    session.execute(query, (
        user_id,
        user.username,
        user.email,
        user.password_hash,
        user.full_name,
        user.bio,
        user.avatar_url,
        created_at
    ))
    return User(
        user_id=user_id,
        username=user.username,
        email=user.email,
        password_hash=user.password_hash,
        full_name=user.full_name,
        bio=user.bio,
        avatar_url=user.avatar_url,
        created_at=created_at
    )

@router.put("/users/{user_id}", response_model=User)
def update_user(user_id: UUID, user: UserCreate):
    # Verifica existencia
    result = session.execute(
        "SELECT * FROM users_by_id WHERE user_id = %s",
        (user_id,)
    ).one()
    if not result:
        raise HTTPException(status_code=404, detail="User not found")
    # Actualiza los campos editables
    session.execute(
        "UPDATE users_by_id SET username=%s, email=%s, password_hash=%s, full_name=%s, bio=%s, avatar_url=%s WHERE user_id=%s",
        (
            user.username,
            user.email,
            user.password_hash,
            user.full_name,
            user.bio,
            user.avatar_url,
            user_id
        )
    )
    # Devuelve el usuario actualizado
    updated = session.execute(
        "SELECT * FROM users_by_id WHERE user_id = %s",
        (user_id,)
    ).one()
    return User(**updated._asdict())

@router.patch("/users/{user_id}", response_model=User)
def partial_update_user(user_id: UUID,
                       username: str = Body(None),
                       email: str = Body(None),
                       password_hash: str = Body(None),
                       full_name: str = Body(None),
                       bio: str = Body(None),
                       avatar_url: str = Body(None)):
    # Verifica existencia
    result = session.execute(
        "SELECT * FROM users_by_id WHERE user_id = %s",
        (user_id,)
    ).one()
    if not result:
        raise HTTPException(status_code=404, detail="User not found")
    # Construye la consulta dinámica solo con los campos enviados
    fields = {}
    if username is not None:
        fields['username'] = username
    if email is not None:
        from pydantic import validate_email
        try:
            validate_email(email)
        except Exception:
            raise HTTPException(status_code=422, detail="Invalid email format")
        fields['email'] = email
    if password_hash is not None:
        fields['password_hash'] = password_hash
    if full_name is not None:
        fields['full_name'] = full_name
    if bio is not None:
        fields['bio'] = bio
    if avatar_url is not None:
        fields['avatar_url'] = avatar_url
    if not fields:
        raise HTTPException(status_code=400, detail="No fields to update")
    set_clause = ', '.join([f"{k}=%s" for k in fields.keys()])
    values = list(fields.values()) + [user_id]
    session.execute(f"UPDATE users_by_id SET {set_clause} WHERE user_id=%s", values)
    # Devuelve el usuario actualizado
    updated = session.execute(
        "SELECT * FROM users_by_id WHERE user_id = %s",
        (user_id,)
    ).one()
    return User(**updated._asdict())

@router.delete("/users/{user_id}", status_code=204)
def delete_user(user_id: UUID):
    # Verifica existencia antes de borrar
    result = session.execute("SELECT user_id FROM users_by_id WHERE user_id = %s", (user_id,)).one()
    if not result:
        raise HTTPException(status_code=404, detail="User not found")
    session.execute("DELETE FROM users_by_id WHERE user_id = %s", (user_id,))
    return
