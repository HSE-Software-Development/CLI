# CLI

[![CI Tests](https://github.com/HSE-Software-Development/CLI/actions/workflows/go.yaml/badge.svg)](https://github.com/HSE-Software-Development/CLI/actions)

Простой интерпретатор командной строки, поддерживающий [самореализованные команды](#поддерживаемые-команды), вызов внешних программ, а также поддержку переменных, своих и окружения.

## Поддерживаемые команды
- cat [FILE] — вывести на экран содержимое файла
- echo — вывести на экран свой аргумент (или аргументы)
- wc [FILE] — вывести количество строк, слов и байт в файле
- pwd — распечатать текущую директорию
- exit — выйти из интерпретатора

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