# Proyecto-Programacion-Avanzada
Entregables

"Utilizamos tres herramientas principales:

Go (Golang) → el lenguaje en el que está escrito el backend.

Gin → un framework que nos ayuda a crear las rutas de la API, como /dolls o /clientes.

GORM → que es un ORM, básicamente una librería que conecta nuestro código con la base de datos sin tener que escribir SQL a mano.

Y como base de datos usamos SQLite, que es una base de datos ligera y sencilla. No necesita servidor, se guarda en un solo archivo."

2. Estructura del proyecto

"Nuestro backend funciona como una API REST. Esto quiere decir que otras aplicaciones, como un frontend o Postman, pueden hablar con él usando peticiones HTTP.

Tenemos tres módulos principales:

Dolls → representan las autómatas.

Clientes → personas que solicitan cartas.

Cartas → los mensajes que redactan las Dolls."

3. Cómo funciona cada CRUD

"Para cada entidad tenemos un CRUD:

GET → listar o ver los datos.

POST → crear nuevos registros.

PUT → actualizar registros existentes.

DELETE → eliminar registros.

Por ejemplo:

Si hago POST /dolls puedo registrar una Doll nueva.

Si hago GET /dolls veo todas las que existen.

Si hago PUT /dolls/1 actualizo la Doll con id 1.

Si hago DELETE /dolls/1 la elimino."

4. Flujo de trabajo (ejemplo)

"Un ejemplo del uso sería:

Primero creo un Cliente (por ejemplo, Gilbert).

Después registro una Doll (como Violet).

Finalmente creo una Carta, que une a ese cliente con esa Doll, indicando la fecha, el contenido y el estado de la carta.

De esta manera, el sistema queda organizado y fácil de consultar."

5. Demostración práctica

"A través de Postman podemos probar la API. Tenemos una colección con todos los endpoints:

Crear, ver, actualizar y borrar Dolls.

Crear, ver, actualizar y borrar Clientes.

Crear, ver, actualizar y borrar Cartas.

Por ejemplo, puedo enviar un POST /cartas con un JSON que diga:

{
  \"cliente_id\": 1,
  \"doll_id\": 1,
  \"fecha\": \"2025-08-17\",
  \"estado\": \"borrador\",
  \"contenido\": \"Querida Violet...\"
}


y automáticamente la carta se guarda en la base de datos."
