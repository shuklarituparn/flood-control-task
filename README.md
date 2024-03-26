# Флудь - Контроль 

---

![Welcome](https://github.com/shuklarituparn/flood-control-task/assets/66947051/e4777094-4d16-4287-bee7-7afe21294d35)

Создал распределенный rate-limiter в go, который поддерживает переменные лимиты скорости.

Сервис доступен здесь: 

 - `https://rate-limiter.rtprnshukla.ru/` (открыт, чтобы увидеть персика)

`https://rate-limiter.rtprnshukla.ru/persik/userID=<int64>&userType=<normal/premium>`

---


### Архитектура

![Architechture](https://github.com/shuklarituparn/flood-control-task/assets/66947051/28438f21-8e4e-4f05-a477-5b614dd30577) 

- Запрос от пользователя проходит через Лод Балансер
- Затем он направляет запрос к одному из трех экземпляров сервиса 
- У каждой сервиса есть "Флуд-Контроль", который считывает глобальный кэш из redis
- Я добавил флуд-контроль в качестве мидлвейра к сервису.

---

### Конфигурация

```yaml
redis:
  address: 'redis:6379'
  password: ''
  db: 0
rate_limiter:
  default_rate_limit:
    rate: 8
    window_seconds: 1s
  user_types:
    normal:
      key_prefix: 'rl:normal'
      rate_limit:
        rate: 5
        window_seconds: 1s
    premium:
      key_prefix: 'rl:premium'
      rate_limit:
        rate: 10
        window_seconds: 1s


```

`https://rate-limiter.rtprnshukla.ru/persik/userID=<int64>&userType=<normal/premium>`

- Eсли локально
`http://localhost:8080/persik/userID=<int64>&userType=<normal/premium>`

- Usertype необходим для определения количества запросов
- По умолчанию пользователю предоставляется rate 8, а затем window 1 секунда
- Сервис также поддерживает обычных пользователей с другим тарифом
- Сервис также поддерживает премиум-пользователей с другим тарифом
---

### Установка

- клонируйте репозиторий, выполнив
   

``` git@github.com:shuklarituparn/flood-control-task.git```

- Выполните следующую команду

    `make setup` чтобы установить Golang и другие зависимости

- убедитесь, что у вас установлен docker

- Выполните следующую команду
```docker compose up```

- вы можете получить доступ к сервису по умолчанию в `http://localhost:8080` и встретиться с персиком

- вы также можете получить доступ к сервису через один из этих портов `http://localhost:8090`, `http://localhost:8091`, `http://localhost:8092`

---

### Тестировать

![Screenshot from 2024-03-25 22-48-09](https://github.com/shuklarituparn/flood-control-task/assets/66947051/44db7a98-a637-4fb9-9431-71267be9518e)
![Screenshot from 2024-03-25 22-48-05](https://github.com/shuklarituparn/flood-control-task/assets/66947051/a9e57922-d07c-4467-a5ad-ecc610afd62b)

`вы также можете использовать curl для тестирования сервиса`

![Screenshot from 2024-03-26 11-10-01](https://github.com/shuklarituparn/Conversion-Microservice/assets/66947051/bbaf8306-e305-491b-a7a4-74a24422aefd)

---

### Что надо было сделать

Реализовать интерфейс с методом для проверки правил флуд-контроля. Если за последние N секунд вызовов метода Check будет больше K, значит, проверка на флуд-контроль не пройдена.

- Интерфейс FloodControl располагается в файле main.go.

- Флуд-контроль может быть запущен на нескольких экземплярах приложения одновременно, поэтому нужно предусмотреть общее хранилище данных. Допустимо использовать любое на ваше усмотрение. 

----

### как я сделал

- Хотел выбрать tarantool, так как я являюсь Амбассадором VK и хотел использовать "наш" продукт (но выбрал redis для более быстрого прототипирования)
- Использовал Redis sorted sets (и sliding window method of rate limiting)
- Использовал nginx в качестве лод балансера
- Хотел сделать флуд-контрол на сайте лод балансера, чтобы получить элегантное решение, но из-за необходимости реализации интерфейса не смог этого сделать 


