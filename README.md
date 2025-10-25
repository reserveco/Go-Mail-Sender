Go-Mail Sender - Полное руководство
===================================

![Flat-square](https://img.shields.io/badge/Go-1.16+-00ADD8?logo=go)  
![Flat-square](https://img.shields.io/badge/SMTP-%D0%BF%D0%BE%D0%B4%D0%B4%D0%B5%D1%80%D0%B6%D0%BA%D0%B0-green)  
![Flat-square](https://img.shields.io/badge/REST-API-blue)  
![Flat-square](https://img.shields.io/badge/%D0%9B%D0%B8%D1%86%D0%B5%D0%BD%D0%B7%D0%B8%D1%8F-MIT-yellow)

Оглавление
----------

*   Обзор программы
    
*   Возможности
    
*   Установка и настройка
    
*   Режим командной строки (CLI)
    
*   REST API режим
    
*   Форматы файлов
    
*   Логирование и мониторинг
    
*   Примеры использования
    
*   Устранение неполадок
    
*   Безопасность
    
*   Производительность
    

Обзор программы
---------------

Go-Mail Sender - это универсальная программа для отправки электронной почты, написанная на Go. Поддерживает два режима работы:

*   **CLI режим** - для разовой отправки писем через командную строку
    
*   **REST API режим** - для интеграции с другими приложениями
    

Возможности
-----------

*   ✅ Отправка одиночных и массовых писем
    
*   ✅ Поддержка HTML и текстовых писем
    
*   ✅ Вложения файлов
    
*   ✅ Загрузка списка получателей из JSON
    
*   ✅ Асинхронная отправка нескольким получателям
    
*   ✅ Детальное логирование в JSON формате
    
*   ✅ REST API для интеграции
    
*   ✅ Поддержка STARTTLS и SSL
    
*   ✅ Автоматическая генерация ID писем
    
*   ✅ Статус отправки для каждого получателя
    

Установка и настройка
---------------------

### Системные требования

*   Ubuntu 18.04+ или другая Linux-система
    
*   Go 1.16+ (для сборки из исходников)
    
*   Доступ к SMTP серверу
    

### Установка Go
```bash
\# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

\# Или установка последней версии
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar \-C /usr/local \-xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' \>> ~/.bashrc
source ~/.bashrc
```
### Сборка программы

```bash
\# Инициализация модуля
go mod init go-mail

\# Установка зависимостей
go get github.com/gin-gonic/gin
go get github.com/google/uuid
go get github.com/wneessen/go-mail

\# Сборка
go build \-o go-mail

\# Проверка
./go-mail \-help
```
### Настройка SMTP

Программа поддерживает различные SMTP серверы. Основные параметры:

Параметр

По умолчанию

Описание

`-host`

`smtp.example.com`

SMTP сервер

`-port`

`25`

Порт SMTP

`-user`

`user@example.com`

Email отправителя

`-pass`

`your_password`

Пароль

Режим командной строки (CLI)
----------------------------

### Все доступные флаги

```bash
./go-mail \-help
```
**Основные флаги:**

Флаг

Описание

Пример

`-host`

SMTP хост

`-host smtp.example.com`

`-port`

SMTP порт

`-port 587`

`-user`

Email отправителя

`-user sender@example.com`

`-pass`

Пароль

`-pass "secret"`

`-to`

Получатели через запятую

`-to "user1@example.com,user2@example.com"`

`-toFile`

JSON файл со списком получателей

`-toFile recipients.json`

`-subject`

Тема письма

`-subject "Важное сообщение"`

`-body`

Текст письма

`-body "Содержимое письма"`

`-bodyFile`

Файл с содержимым письма

`-bodyFile message.html`

`-html`

Использовать HTML формат

`-html`

`-attach`

Файлы для вложения через запятую

`-attach "file1.pdf,image.jpg"`

`-id`

ID письма

`-id "newsletter_001"`

`-api`

Запуск в API режиме

`-api :8080`

### Примеры использования CLI

#### 1\. Простое текстовое письмо

bash

./go-mail \\
  \-to "user@example.com" \\
  \-subject "Тестовое письмо" \\
  \-body "Это тестовое сообщение"

#### 2\. HTML письмо с файлом

bash

./go-mail \\
  \-to "user@example.com" \\
  \-subject "HTML письмо" \\
  \-bodyFile "newsletter.html" \\
  \-html

#### 3\. Письмо с вложениями

bash

./go-mail \\
  \-to "user@example.com" \\
  \-subject "Документы" \\
  \-body "Смотрите вложения" \\
  \-attach "document.pdf,image.jpg"

#### 4\. Массовая рассылка

bash

./go-mail \\
  \-to "user1@example.com,user2@example.com,user3@example.com" \\
  \-subject "Массовая рассылка" \\
  \-body "Общее сообщение для всех"

#### 5\. Использование файла получателей

```bash
./go-mail \\
  \-toFile "recipients.json" \\
  \-subject "Рассылка из файла" \\
  \-body "Сообщение для списка получателей"
```
#### 6\. Полный пример со всеми параметрами

```bash
./go-mail \\
  \-host "smtp.example.com" \\
  \-port 25 \\
  \-user "sender@example.com" \\
  \-pass "password" \\
  \-to "user1@example.com,user2@example.com" \\
  \-subject "Важное сообщение" \\
  \-bodyFile "message.html" \\
  \-html \\
  \-attach "report.pdf" \\
  \-id "newsletter\_001"
```
REST API режим
--------------

### Запуск API сервера

```bash
./go-mail \-api :8080
```
Сервер запустится на указанном порту и будет доступен для HTTP запросов.

### Эндпоинты API

#### 1\. Отправка email

**POST** `/api/send`

**Тело запроса:**

```json
{
  "to": \["user1@example.com", "user2@example.com"\],
  "subject": "Тема письма",
  "body": "Текст письма",
  "body\_file": "путь/к/файлу.html",
  "is\_html": true,
  "attachments": \["file1.pdf", "image.jpg"\],
  "message\_id": "custom\_id\_123"
}
```
**Ответ:**

```json
{
  "id": "custom\_id\_123",
  "status": "accepted",
  "message": "Письмо принято в обработку"
}
```

#### 2\. Проверка статуса

**GET** `/api/status/{id}`

**Ответ:**

```json
{
  "id": "custom\_id\_123",
  "timestamp": "2024-01-15T14:30:45.123456789Z",
  "overall\_status": "partial\_success",
  "details": \[
    {
      "email": "user1@example.com",
      "status": "success",
      "time": "2024-01-15T14:30:46.123456789Z",
      "message": "Письмо успешно отправлено"
    },
    {
      "email": "user2@example.com",
      "status": "error",
      "time": "2024-01-15T14:30:47.123456789Z",
      "error": "ошибка отправки: ..."
    }
  \]
}
```
#### 3\. Проверка здоровья

**GET** `/api/health`

**Ответ:**

```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T14:30:45.123456789Z",
  "service": "Email Sender API"
}
```
### Примеры использования API

#### Отправка письма через curl

```bash
curl \-X POST http://localhost:8080/api/send \\
  \-H "Content-Type: application/json" \\
  \-d '{
    "to": \["test@example.com"\],
    "subject": "Тест через API",
    "body": "Это тестовое письмо отправленное через REST API",
    "is\_html": false
  }'
```
#### Проверка статуса

```bash
curl http://localhost:8080/api/status/custom\_id\_123
```
#### Использование с Python

```python
import requests

response \= requests.post("http://localhost:8080/api/send", json\={
    "to": \["user@example.com"\],
    "subject": "Письмо из Python",
    "body": "Сообщение отправленное через Python скрипт",
    "is\_html": False
})

print(response.json())
```
Форматы файлов
--------------

### 1\. Файл получателей (JSON)

**Формат:**

```json
{
  "emails": \[
    "user1@example.com",
    "user2@example.com",
    "user3@example.com",
    "user4@example.com"
  \]
}
```
**Пример создания:**

```bash
cat \> recipients.json << EOF
{
  "emails": \[
    "user1@example.com",
    "user2@example.com",
    "user3@example.com"
  \]
}
EOF
```
### 2\. HTML шаблон письма

**Пример:**

```html
<!DOCTYPE html\>
<html\>
<head\>
    <meta charset\="UTF-8"\>
    <title\>Новости компании</title\>
    <style\>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .header { background: #f4f4f4; padding: 20px; text-align: center; }
        .content { padding: 20px; }
        .footer { background: #333; color: white; padding: 10px; text-align: center; }
    </style\>
</head\>
<body\>
    <div class\="header"\>
        <h1\>Новости нашей компании</h1\>
    </div\>
    
    <div class\="content"\>
        <h2\>Уважаемый клиент!</h2\>
        <p\>Мы рады сообщить вам о последних новостях:</p\>
        
        <ul\>
            <li\><strong\>Новый продукт</strong\> - уже доступен</li\>
            <li\><strong\>Специальные предложения</strong\> - только этой неделей</li\>
            <li\><strong\>Обновление сервиса</strong\> - улучшенная производительность</li\>
        </ul\>
        
        <p\>С уважением,<br\>Команда компании</p\>
    </div\>
    
    <div class\="footer"\>
        <p\>&copy; 2024 Наша компания. Все права защищены.</p\>
        <p\><a href\="#" style\="color: #fff;"\>Отписаться от рассылки</a\></p\>
    </div\>
</body\>
</html\>
```
Логирование и мониторинг
------------------------

### Структура логов

Программа создает папку `logs` и сохраняет в нее JSON файлы:

```text
logs/
├── newsletter\_001.json
├── 550e8400-e29b-41d4-a716-446655440000.json
└── custom\_id\_123.json
```
### Формат лог-файла

```json
{
  "id": "newsletter\_001",
  "timestamp": "2024-01-15T14:30:45.123456789Z",
  "from": "sender@example.com",
  "subject": "Новости компании",
  "overall\_status": "partial\_success",
  "recipients": \[
    {
      "email": "user1@example.com",
      "status": "success",
      "time": "2024-01-15T14:30:46.123456789Z",
      "message": "Письмо успешно отправлено"
    },
    {
      "email": "user2@example.com",
      "status": "error",
      "time": "2024-01-15T14:30:47.123456789Z",
      "error": "ошибка отправки: dial tcp: lookup example.com: no such host"
    }
  \]
}
```
### Статусы отправки

*   `success` - все письма доставлены
    
*   `partial_success` - часть писем доставлена, часть нет
    
*   `failed` - все письма не доставлены
    
*   `accepted` - письмо принято в обработку (только для API)
    

### Мониторинг логов

```bash
\# Просмотр всех логов
ls \-la logs/

\# Чтение конкретного лога
cat logs/newsletter\_001.json

\# Мониторинг в реальном времени (если используется API)
tail \-f logs/latest.log
```
Примеры использования
---------------------

### Сценарий 1: Ежедневная рассылка

```bash
#!/bin/bash
\# daily\_newsletter.sh

./go-mail \\
  \-toFile "subscribers.json" \\
  \-subject "Ежедневные новости $(date +%Y-%m-%d)" \\
  \-bodyFile "daily\_news.html" \\
  \-html \\
  \-id "daily\_$(date +%Y%m%d)"
```
### Сценарий 2: Уведомления системы

```bash
#!/bin/bash
\# system\_alert.sh

./go-mail \\
  \-to "admin@company.com" \\
  \-subject "Системное оповещение" \\
  \-body "Система обнаружила проблему: $1" \\
  \-id "alert\_$(date +%s)"
```
### Сценарий 3: Интеграция с веб-приложением

```bash
\# Запуск API сервера как службы
nohup ./go-mail \-api :8080 \> mailer.log 2\>&1 &

\# Использование в скрипте
curl \-X POST http://localhost:8080/api/send \\
  \-H "Content-Type: application/json" \\
  \-d '{
    "to": \["'"$USER\_EMAIL"'"\],
    "subject": "Регистрация завершена",
    "body": "Добро пожаловать в нашу систему!",
    "message\_id": "reg\_'$(date +%s)'"
  }'
```
Устранение неполадок
--------------------

### Частые проблемы и решения

#### 1\. Ошибки подключения к SMTP

```text
Ошибка подключения: dial tcp: lookup host: no such host
```
**Решение:**

*   Проверьте правильность имени SMTP сервера
    
*   Убедитесь в доступности сети
    
*   Проверьте firewall настройки
    

#### 2\. Ошибки аутентификации

```text
Ошибка аутентификации: 535 5.7.8 Error: authentication failed
```
**Решение:**

*   Проверьте логин и пароль
    
*   Для Gmail используйте "Пароль приложения"
    
*   Убедитесь что аккаунт не заблокирован
    

#### 3\. Письма не доставляются

**Решение:**

*   Проверьте настройки SPF/DKIM домена
    
*   Убедитесь что IP адрес не в черном списке
    
*   Проверьте папку "Спам" получателя
    

#### 4\. Ошибки с вложениями

```text
Ошибка чтения файла: open file.pdf: no such file or directory
```
**Решение:**

*   Убедитесь что файл существует
    
*   Проверьте права доступа к файлу
    
*   Используйте полный путь к файлу
    

### Диагностика

#### Проверка SMTP соединения

```bash
telnet smtp.example.com 25
```
#### Проверка DNS записей

```bash
\# Проверка MX записей
dig MX example.com

\# Проверка SPF
dig TXT example.com
```
#### Тестовые команды

```bash
\# Простой тест отправки
./go-mail \-to "test@example.com" \-subject "Test" \-body "Test message"

\# Тест с подробным выводом
./go-mail \-to "test@example.com" \-subject "Debug" \-body "Test" \-id "debug\_test"
```
### Логи для отладки

Программа выводит подробную информацию в консоль:

```text
📧 Отправка письма ID: newsletter\_001
   Сервер: smtp.example.com:25
   От: sender@example.com
   Тема: Тестовое письмо
   Получатели: user1@example.com, user2@example.com
🚀 Асинхронная отправка 2 писем...
   ✅ Успешно: user1@example.com
   ❌ Ошибка: user2@example.com - ошибка отправки: ...
✅ Процесс отправки завершен! Статус: partial\_success
📁 Лог сохранен: logs/newsletter\_001.json
```
Безопасность
------------

### Рекомендации по безопасности

1.  **Хранение паролей:**
    
    *   Не храните пароли в скриптах
        
    *   Используйте переменные окружения
        
    *   Рассмотрите использование секретов
        
2.  **Доступ к API:**
    
    *   Запускайте API только на localhost в production
        
    *   Используйте reverse proxy с аутентификацией
        
    *   Ограничьте доступ по IP
        
3.  **Валидация ввода:**
    
    *   Проверяйте email адреса на валидность
        
    *   Ограничивайте размер вложений
        
    *   Санкционируйте имена файлов
        

### Пример безопасной конфигурации

```bash
\# Использование переменных окружения
export SMTP\_PASSWORD\="your\_password"
./go-mail \-to "$RECIPIENT" \-pass "$SMTP\_PASSWORD"

\# Запуск API только на localhost
./go-mail \-api 127.0.0.1:8080
```


Производительность
------------------

### Оптимизация для массовых рассылок
1.  **Лимиты отправки:**
    *   Настройте задержки между отправками
    *   Разбивайте большие списки на части
    *   Используйте очередь отправки
        
2.  **Мониторинг ресурсов:**
    *   Следите за использованием памяти
    *   Мониторьте сетевую активность
    *   Логируйте время выполнения
        

### Пример оптимизированного скрипта

```bash
#!/bin/bash
\# bulk\_sender.sh

BATCH\_SIZE\=50
DELAY\=5

\# Разбиваем большой список на части
jq \-c '.emails | .\[\]' recipients.json | \\
while read email; do
    ./go-mail \\
        \-to "$email" \\
        \-subject "Массовая рассылка" \\
        \-bodyFile "message.html" \\
        \-html \\
        \-id "bulk\_$(date +%s)"
    
    \# Задержка между отправками
    sleep $DELAY
done
```
Лицензия
--------
Этот проект распространяется под лицензией MIT. Подробности см. в файле LICENSE.


Поддержка
---------
Если у вас возникли вопросы или проблемы:
1.  Проверьте документацию и примеры использования
2.  Изучите логи в папке `logs/`
3.  Убедитесь в правильности настроек SMTP
    

* * *

**Полезные ссылки:**
*   [Документация Go](https://golang.org/doc/)
*   [Библиотека go-mail](https://github.com/wneessen/go-mail)
*   [Фреймворк Gin](https://gin-gonic.com/)
    


Удачи в использовании! 🚀
