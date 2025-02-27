-- +goose Up
CREATE TABLE IF NOT EXISTS "apps" ("id" text NOT NULL, "installedAt" timestamp NOT NULL);
