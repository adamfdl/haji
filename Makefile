build:
	cd $(PWD)/cmd && go build -v -o tenx

run:
	JWT_KEY=randomstring \
	DB_HOST=salt.db.elephantsql.com \
	DB_PORT=5432 \
	DB_USER=vpfwajkv \
	DB_NAME=vpfwajkv \
	DB_PASSWORD=H1worPYJd4DiGN7zD6UlxpMHzmYU1-Zy \
	cmd/tenx