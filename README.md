# 🔍 Sherlock - Telegram Channel Monitor

Мониторинг подписчиков вашего Telegram-канала с мгновенными уведомлениями о подписках и отписках.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Telegram](https://img.shields.io/badge/Telegram-Bot-2CA5E0?logo=telegram)
![SQLite](https://img.shields.io/badge/SQLite-Database-003B57?logo=sqlite)

## ✨ Возможности

- 🔔 **Мгновенные уведомления** о новых подписчиках и отписках
- 📊 **Статистика** подписчиков прямо в Telegram
- 💾 **Автосохранение** данных в SQLite базу
- ⏰ **Отслеживание длительности** подписки
- 🛡️ **Защита данных** - информация сохраняется даже при отписке

## 🚀 Быстрый старт

### 1. Клонирование и настройка

```bash
git clone https://github.com/your-username/sherlock.git
cd sherlock
cp configs/config.example.yml configs/config.yml
```
```
2. Создание бота через @BotFather
Найдите @BotFather в Telegram

Отправьте /newbot

Следуйте инструкциям и получите токен

Дайте боту права администратора в вашем канале

3. Настройка конфигурации
configs/config.yml:

yaml
bot:
  token: "1234567890:ABCDEFGHIJKLMNOPQRSTUVWXYZ"  # Токен от @BotFather

channel:
  id: -1001234567890                              # ID вашего канала

admin:
  chat_id: 123456789                              # Ваш Chat ID

database:
  path: "data/sherlock.db"                        # Путь к базе данных

monitor:
  interval: 2                                     # Интервал проверки (минуты)
Или через .env:

env
BOT_TOKEN=1234567890:ABCDEFGHIJKLMNOPQRSTUVWXYZ
CHANNEL_ID=-1001234567890
ADMIN_CHAT_ID=123456789
4. Запуск
bash
# Установка зависимостей
go mod download

# Запуск приложения
go run cmd/main.go
📋 Получение необходимых ID
🔍 Получение ID канала
Добавьте @username_to_id_bot в канал

Перешлите любое сообщение из канала боту

Получите ID (начинается с -100)

👤 Получение вашего Chat ID
Напишите @userinfobot в Telegram

Бот ответит с вашим Chat ID

🏗️ Структура проекта
text
sherlock/
├── cmd/
│   └── main.go                 # Точка входа
├── configs/
│   └── config.yml              # Конфигурация
├── data/
│   └── sherlock.db             # База данных (автосоздается)
├── pkg/
│   ├── channel/                # Логика мониторинга канала
│   │   ├── monitor.go          # Основной монитор
│   │   ├── database.go         # Работа с БД
│   │   ├── notifications.go    # Уведомления
│   │   └── types.go            # Структуры данных
│   └── repository/
│       └── sqlite.go           # Инициализация SQLite
└── .env                        # Переменные окружения
🎯 Использование
Команды бота:
/stats - Показать статистику подписчиков

/help - Показать справку

/start - Запустить мониторинг

Пример уведомлений:
text
🎉 Новый подписчик!

👤 Имя: Иван Иванов
📎 Username: @ivanov
🆔 User ID: 123456789
⏰ Время: 20:00:05
text
😢 Кто-то отписался

👤 Имя: Петр Петров  
📎 Username: @petrov
🆔 User ID: 987654321
⏰ Был с нами: 2 дня 3 часа 15 минут
⚙️ Настройка
Интервал проверки
yaml
monitor:
  interval: 5  # минуты (рекомендуется 2-5)
База данных
yaml
database:
  path: "data/sherlock.db"  # Можно изменить путь
🐛 Решение проблем
Бот не видит подписчиков
bash
# Проверьте права бота в канале
# Бот должен быть администратором
# Должны быть включены права:
# - "Добавлять участников"
# - "Просматривать сообщения"
Ошибка "Not Found"
bash
# Проверьте правильность токена
curl "https://api.telegram.org/botYOUR_TOKEN/getMe"
Просмотр базы данных
bash
sqlite3 data/sherlock.db
sqlite> .tables
sqlite> SELECT * FROM subscribers;
```
