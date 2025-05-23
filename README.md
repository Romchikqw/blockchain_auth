# Прототип децентрализованной системы авторизации

Данный проект представляет собой прототип системы авторизации, основанный на технологии распределённого реестра с использованием Hyperledger Fabric и модели контроля доступа ABAC.

## О проекте
Разработан прототип системы, включающий:
- **Сервер аутентификации** — проверка логина/пароля и выдача JWT токенов.
- **Сервер проверки прав доступа** — проверка валидности токена и обращение к смарт-контракту в блокчейне.
- **Сервер ресурса** — предоставление доступа к защищённым данным после валидации токена через блокчейн.
- **Смарт-контракты** (Chaincode) для хранения атрибутов пользователей, политик доступа и проверки токенов в сети Hyperledger Fabric.

## Технологии
- Hyperledger Fabric
- Go (Golang)
- PostgreSQL
- Docker / Docker Compose
- JWT (JSON Web Tokens)
- TLS-шифрование

## Структура проекта
- `auth-server/` — сервер аутентификации пользователей.
- `abac-server/` — сервер проверки прав доступа.
- `resource-server/` — сервер защищённого ресурса.
- `abac-chaincode/` — смарт-контракты для Hyperledger Fabric.
- `fabric-samples-main/` — тестовая сеть Hyperledger Fabric для развертывания блокчейна.

## Основные возможности
- Аутентификация пользователя и выдача токена доступа.
- Проверка прав доступа с использованием атрибутов и политик.
- Хранение информации о пользователях, правах и токенах в неизменяемом реестре (блокчейн).
- Защита от фальсификации токенов и атак на целостность данных.

## Цель
Создание децентрализованной системы авторизации, устойчивой к атакующим воздействиям, обеспечивающей прозрачность процессов и неизменность данных с помощью технологии блокчейна.

---

> Данный проект был разработан в рамках научно-исследовательской работы по теме "Прототип децентрализованной системы авторизации".

