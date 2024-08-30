CREATE TABLE "todos" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "title" character varying NOT NULL,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
)