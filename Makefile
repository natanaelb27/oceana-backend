db-init:
	psql -U postgres -h localhost -d postgres -f db_init.sql