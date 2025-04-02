# Contribution

Для того, чтобы добавить обработку новых команд в CLI, вам необходимо:
- Создать fork этого проекта
- В файле [commands.go](../internal/executor/commands.go) добавить функцию со следующей сигнатурой:
```go
func name_cmd(cmd parseline.Command, b *bytes.Buffer) (*bytes.Buffer, error) {}
```
- Напиши тесты для новой функции в фаил [commands_test.go](../internal/executor/commands_test.go)
```go
//Напишите название новой функции
func TestNameCommand(t *testing.T) {
	tests := []struct {
		name    string // название теста
		cmd     parseline.Command // Команда, которую хотите выполнить
		input   *bytes.Buffer // буффер от прошлой функции в pipeline
		want    string // Ожидаемый результат
		wantErr bool // Должена ли функция возвращать ошибку
	}{
        // Задайте параметры тестов
	}
    for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            // Вместо NAME_CMD напишите название новой функции
			got, err := NAME_CMD(tt.cmd, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("wc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.String() != tt.want {
				t.Errorf("wc() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}
```

- После прохождения тестов, добавьте новый функционал в систему. Для этого вернемся в фаил [commands.go](../internal/executor/commands.go) и в методе New() добавьте инициализацию вашей команды.
```go
func newCommands() commands {
	cmds := make(commands)

	// Here you can add new command in CLI
	// cmd["name_command"] = name_command
	// below you need to implement a command with the following signature:
	// func(parseline.Command, *bytes.Buffer) (*bytes.Buffer, error)
	cmds["cat"] = cat
	cmds["echo"] = echo
	cmds["exit"] = exit
	cmds["pwd"] = pwd
	cmds["wc"] = wc
	cmds["grep"] = grep


	return cmds
}
```

Требования/советы к реализации:
- На вход подается структура parseline.Command, которая задана в файле [parser.go](../internal/parseline/parser.go) 
```go
// Command store name of command and it's flags and args
type Command struct {
	Name string  
	Args []string 
}
```

- Парсер разделяет получаемую команду на command и args, в буффере будет записана только возвращаемое значение предыдущей команды в pipeline. Буффер всегда будет инициализован.

- В случае если функция завершилась без ошибки, но ничего не возвращает, возвращать не nil буффер!!!
```go 
return bytes.NewBufferString(""), nil // так


return nil, nil // не так
```
- В случае если во время испольнения произошла какая-то ошибка, вы вольны создавать ее как через пакет errors, так и через пакет fmt
```go
return nil, errors.New("no input provided") // так, в случае создания ошибки

return nil, fmt.Errorf("cat: %w", err) // так, в случае обертки надо полученной ошибкой
```
- Parser соберет command без удаление скобок и ковычек, примеры:
```
>>> echo "111"
```
На испольнение будет передана команда:
```go
Command {
    name: "echo",
    args: [`"111"`],
}
```


