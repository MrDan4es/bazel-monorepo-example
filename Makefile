deps-update: ## Обновляет все go зависимости для bazel
	go mod tidy && bazel run :gazelle-update-repos && bazel run :gazelle

bundle:
	bazel run :bundle ## Собрать все контейнеры

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m- %-14s\033[0m \033[32m%s\033[0m\n", $$1, $$2}'

.DEFAULT_GOAL := help
