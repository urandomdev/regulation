package ent

//go:generate atlas migrate diff baseline --dir "file://migrate/migrations" --to "ent://schema" --dev-url "docker+postgres://ghcr.io/freezm-ltd/postgres:latest/postgres?search_path=public" --format "{{sql . \"  \"}}"
