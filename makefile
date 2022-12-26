# SHELL := /bin/zsh

run:
	$(shell grep -v '#' .env | tr "\n" " ") go run .
