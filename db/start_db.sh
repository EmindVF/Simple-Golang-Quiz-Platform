# Build postgres image
docker build -t postgres_db . 

# Run the database container
docker run -d -p 5432:5432 --name quiz_platform postgres_db

# Pull pgAdmin container
#docker pull dpage/pgadmin4

# Run the pgAdmin container
#docker run -p 9090:80 \
#--name pgadmin \
#-e 'PGADMIN_DEFAULT_EMAIL=test@gmail.com' \
#-e 'PGADMIN_DEFAULT_PASSWORD=password' \
#-d dpage/pgadmin4