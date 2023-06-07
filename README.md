# test-task-ozon
Запуск с in-memory хранилищем:<br>
METHOD=in-memory docker-compose up --build<br>
Запуск с PostgreSQL:<br>
docker-compose up --build <br>
(можно и так: METHOD=postgres docker-compose up --build)<br>
Для создания новой, короткой ссылки используется POST-запрос по адресу localhost:8080/api/links.<br>
Тело запроса должно быть представлено в JSON-формате, следующего вида:<br>
{"initial_url":"https://www.ozon.ru/travel/?mwc_campaign=oztravel_horizontal-menu_flight"}<br>
Для получения оригинальной ссылки необходимо совершить GET-запрос по адресу, который вернулся POST-запросом.<br>
Короткий адрес имеет подобный вид:<br>
localhost:8080/KrtG99MTpQ
