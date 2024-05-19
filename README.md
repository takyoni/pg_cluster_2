# Практическое изучение отказоустойчивости, часть 2. Верификация системы обработки отказа СУБД

Тесты проводились на модифицированном кластере из задания 1: в кластер была добавлена сущность writer, необходимая для записи данных в СУБД, pgadmin для работы с кластером.

# Результаты тестирования

Тестируемый случай: на мастере происходит какая-то ошибка и он становится недоступен (ошибка моделируется при помощи iptables), на slave происходит promote до master и данные начинают поступать напрямую в slave.

При ```synchronous_commit = off``` произошла потеря 42 записей из 1 миллиона.

При ```synchronous_commit = remote_apply``` потери данных не обнаружено.

# Запуск

```
docker compose up
```
