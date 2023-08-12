h1 Тестовое задание для соискателя на вакансию
backend-разработчик.

1. Тестовое задание необходимо выполнять в Docker инфраструктуре.
2. В контейнере необходимо развернуть MySQL (MariaDB)
3. Так же необходимо развернуть rabbitMQ, который будет состоять из одного
   консьюмера и одного продьюссера
4. Консьюмер и продьюссер должны быть написаны на “голом” PHP либо Golang.
5. Необходимо применить супервизор для отслеживания работоспособности
   консьюмера и запустить консьюмер в 5 потоков.
   Продьюсер должен отправлять в сообщении URL-адрес. Кол-во сообщений и
   конкретные URL адреса могут быть произвольные.
   Консьюмер должен поглощать сообщения с задержкой в 30 секунд и выполнять
   запрос по указанному адресу. Консьюмер должен записывать в таблицу БД код ответа,
   response header, и контент ответа. (структура таблицы произвольная).
   Если код ответа не 200, необходимо повторно один раз отослать сообщение
   через rabbitMQ с задержкой в 15 секунд.
   Создать запрос на выборку, в котором будет выводиться общее кол-во запросов
   и кол-во запросов в response header которых встречается поле 'new' со значением 1
   (название поля может быть любым).