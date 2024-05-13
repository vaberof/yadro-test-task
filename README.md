# Использование

## Запуск докер контейнера с приложением

При создании контейнера в него копируются исполняемый файл *main* и тестовые файлы из папки *examples*.
В качестве аргумента **TEST_FILE** нужно указывать название тестового файла из папки *examples*

    make docker.run TEST_FILE=test_file_ok_1.txt

## Локальный запуск

    make build.linux
    ./cmd/yadro-test-task/build/main examples/test_file_ok_1.txt

## Запуск юнит-тестов

    make tests.run
