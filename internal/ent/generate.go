package ent

//go:generate go tool ent generate ./schema --feature sql/lock,sql/upsert,schema/snapshot,sql/modifier,sql/execquery --template generate_extension.go.tmpl
