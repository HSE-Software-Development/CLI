# CLI

[![CI Tests](https://github.com/HSE-Software-Development/CLI/actions/workflows/go.yaml/badge.svg)](https://github.com/HSE-Software-Development/CLI/actions)

Простой интерпретатор командной строки, поддерживающий [самореализованные команды](#поддерживаемые-команды), вызов внешних программ, а также поддержку переменных, своих и окружения.

## Поддерживаемые команды
Данный интерпретатор уже поддерживает многие команды, такие как cat, echo, wc, grep и многие другие.
Подробнее о них можете узнать в [Commands.md](docs/Commands.md)

## Запуск под macOS/Linux
``` bash
# build
chmod +x scripts/build.sh
./scripts/build.sh

# Run
chmod +x scripts/run.sh
./scripts/run.sh
```

## Запуск под Windows
``` bash
# build
scripts\build.bat

#Run
#scripts\run.bat
.\bin\cli-app.exe
```

## Переменные окружения

При запуске, программа подгружает переменные окружения с вашего устройства. 
При запуске под unix-подобной системой это будут:
- "PWD", "SHELL", "TERM", "USER", "OLDPWD", "LS_COLORS", "MAIL", "PATH", "LANG", "HOME", "_*"

## Тестирование

Для запуска всех unit-test программы, запустите:
```bash
go test -v ./...
```

## Contribution

Если вам так понравился наш продукт, что у вас появилось желание его доработать, вы всегда можете добавить в него свой функционал. 

[Как мне это сделать?](docs/Contribution.md)

