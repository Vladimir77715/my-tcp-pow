# Реализация защиты от DDOS атак при помощи Proof of Work

Данный репозиторий содержит реализацию механизма защиты от DDOS атак с использованием Proof of Work.

Каждый запрос от клиента перед доступом к целевому ресурсу должен пройти проверку, соответствующую алгоритму Proof of Work.

В данной реализации выбран алгоритм "Вычисление хеш-функции с условием".
Клиенту предлагается найти такое значение (часто называемое "nonce"), которое, будучи добавленным к определённым данным (например, к содержимому запроса), после применения хеш-функции (например, SHA-256) даст результат, начинающийся с определённого количества ведущих нулей. 

Этот алгоритм был выбран из-за его оптимального соотношения нагрузки, которое не позволит слишком слабым устройствам справиться с ним, но при этом затруднит DDOS атаки.

## Установка mockery

```bash
go install github.com/vektra/mockery/v2@v2.42.0
```

## Сборка клиента 

```bash
 docker build -f client.Dockerfile -t Тег:Версия
```

## Сборка Сервера  

```bash
 docker build -f server.Dockerfile -t Тег:Версия
```
## Переменные окражения 

###Сервер
```bash
SERVER_ADDRESS - адресс tpc сервера
SERVER_PORT - порт tpc сервера
SERVER_MIN_SOLUTION_RANGE - минимальное значение, которые нужно получить в хеш функции 
SERVER_MAX_SOLUTION_RANGE - максимальное значение, которые нужно получить в хеш функции
```

###Клиент
```bash
SERVER_ADDRESS - адресс tpc сервера c портом
```
