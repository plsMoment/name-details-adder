CREATE TYPE "genders" AS ENUM ('male', 'female');
CREATE TABLE "users" (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "surname" VARCHAR NOT NULL,
    "patronymic" VARCHAR,
    "age" INTEGER NOT NULL,
    "gender" genders NOT NULL,
    "nationality" VARCHAR NOT NULL
);