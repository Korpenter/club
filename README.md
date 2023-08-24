# YADRO test task

## Инструкции по запуску в Docker
Собрать образ:

```shell
sudo docker buildx build -t club .  
```

Запустить контейнер и указать обрабатываемый файл:

```shell
docker run -v /path/to/user/input.txt:/data/input.txt club
```
