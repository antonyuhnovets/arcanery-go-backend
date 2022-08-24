# Arcanery: go-backend

<h3>Packages and structure</h3>

    cmd/daemon:
      - main.go

    internal:
      - app: setup and load server
      - controller/http/v1 - router and handlers, gin framework for https
      - domain/models - structs with main entities
      - service - usecases for repository, room service, abstract interfaces (interactors)

    pkg:
      - websocket - client, room and hub on ws
      - httpserver - server with options  
      - uuid - random id gen

    config: 
      - load env vars, set default host port etc.

    doc: 
      - swagger documentation


<h3>Routes</h3>

   	 v1/room/..:
	      /new - create new room with random id
          (ws)/$roomId - connect to room by id using websocket
          /rm/$roomId - shutdown and remove room
