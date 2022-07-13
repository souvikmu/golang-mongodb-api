# Go Lang server backend API with MongoDB database integration

### This is a simple backend `API` using `Golang` with `MongoDB` used as a database sytem for Models, with valid `CRUD` operations

</br></br>

> ## API endpoints

</br>

| Endpoints         |     | Method   |     |     |     | Operation                    |
| ----------------- | --- | -------- | --- | --- | --- | ---------------------------- |
| `/api/movies`     |     | `GET`    |     |     |     | Gets all the movies          |
| `/api/movie`      |     | `POST`   |     |     |     | Add a movie to the db        |
| `/api/movie/{id}` |     | `PUT`    |     |     |     | Update a singular movie info |
| `/api/movie/{id}` |     | `DELETE` |     |     |     | Delete a particular movie    |
| `/api/deleteall`  |     | `DELETE` |     |     |     | Erase the entire db          |
