# Добавление приложения как сервис системы с инициализацией переменных окружения

1. Разместите конфигурационный файл сервиса (`user-service.conf`) в директории `/etc/systemd/system`

2. Создайте директорию `/etc/systemd/system/user-service.service.d` и разместите файл с переменными окружения для приложения (`override.conf`)

3. Проверьте работоспособность конфигурации используя команду `systemctl start user-service`

4. Добавьте сервис в автозапуск используя команду `systemctl enable user-service`