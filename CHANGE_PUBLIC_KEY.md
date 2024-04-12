# Процесс обновления публичного ключа

# 1. Генерация сообщения

Администратор генерирует сообщение для подписи на linux

Параметры:1. Публичные ключи всех валидаторов которые были указаны при установке/инициализации acl.

Для обновления публичного ключа обязательно подписывать 3мя валидаторами.



```shell
export validatorPublicKeys="A4JdE9iZRzU9NEiVDNxYKKWymHeBxHR7mA8AetFrg8m4,\
5Tevazf8xxwyyKGku4VCCSVMDN56mU3mm2WsnENk1zv5,\
6qFz88dv2R8sXmyzWPjvzN6jafv7t1kNUHztYKjH1Rd4"
```

1. Адрес по которому меняем публичный ключ, формат base58 сheck:
```shell
export changedAddr="2GFkmC1RE1kMe1HcdcW9Sk3d7nBgNtxeDPRXK7xxxrXordZa7b"
```
1. Новый публичный ключ, формат base58:
```shell
export newPkey="BREP5CVURcJ6CoTdUpJdSNgZThMXvubFueyviSHDRW4Y"
```
Пример запуска команды для ACL 0.2.0:
```shell
./cli generateMessage acl changePublicKey $validatorPublicKeys $changedAddr $newPkey
```

Пример запуска команды для ACL 0.3.1:

**reason** - текстовое поле комментария. Например: "На основании исполнительного листа №бла-
бла-бла от такого-то числа"
**reasonId** - целочисленное поле, которое на бэке будет определяться enum'ом (например: 0-
передача управления кошельком в силу закона, 1-возврат управления пользователю, 2-утрата
публичного ключа пользователя)

```shell
export reason="lost_key"
export reasonId=2
./cli generateMessage acl changePublicKey $validatorPublicKeys $changedAddr $reason $reasonId $newPkey
``` 

Результат:
В директории рядом с приложением "cli" появился файл message.txt

Архив для валидатора должен содержать:
- message.txt
- cli-windows-amd64.exe
- cli-windows-386.exe
  Создать и передать zip архив Валидатору.

# 2. Подпись сообщения

  Валидатор должен запустить командную строку на Windows:
1. Два раза кликните по файлу run_cmd.bat
2. Замените SECRET_KEY_PUT_HERE на приватный ключ в формате base58 или hex.
   ШАБЛОН КОМАНДЫ:

```shell
cli-windows-amd64.exe -s SECRET_KEY_PUT_HERE signMessage
```

Пример команды, которая должна получиться в итоге:
```shell
cli-windows-amd64.exe -s 3aDebSkgXq37VPrzThboaV8oMMbYXrRAt7hnGrod4PNMnGfXjh14TY7cQs8eVT46C4RK4ZyNKLrBmyD5CYZiFmkr signMessage
```
Внимание в зависимости от системы windows нужно запустить либо cli-windows-amd64.exe либо cli-windows-386.exe

3. Вводим команду в открытом конке (после шага 1) и нажимаем ENTER для выполнения команды.

4. В результате выполенения команды в директории рядом с приложением "cli" появился файл signature-*.txt.
В названии файла вместо * указан публичный ключ того кто подписал сообщение.
Скопируйте этот файл и отправьте Администратору.

# 3. Подготовка к обновлению публичного ключа в hlf
3.1. Администратор собирает файлы signature-*.txt у валидаторов и кладет их в директорию
рядом с приложением "cli".
3.2. На этом шаге необходимо проверить, что рядом с приложением cli лежит сообщение
message.txt
3.3. Администратр настраивает файл config_test.yaml, а также сохраняет крипто материалы для
hlf в папку crypto в соответствии с описанием в файле конфигурации config_test.yaml
3.3.1. Перед началом работы нужно сохранить файл config_test_tmp.yaml с названием
config_test.yaml. Файл config_test.yaml используется по умолчанию при вызове cli, как файл
конфигурации для hlf. Если вы хотите использывать файл с другим названием, то его можно
задать через параметр `./cli --cfg config_test.yaml ....` или `./cli -f config_test.yaml ....`
3.3.2. Изменить файл конфигурации для hlf config_test.yaml.
3.3.3. Название канала и полиси. В файле нужно задать название канала куда мы делаем запрос,
а также список пиров в соответствии с полиси для вашего канала.В моем случае название канала acl. 22 строка в файле config_test_tmp.yaml. **Внимание домены
должны соответствовать вашему стенду.

```yaml
acl:
    peers:
        peer0.testnet.uat.dlt.testnet.ch: { }
        peer1.testnet.uat.dlt.testnet.ch: { }
        peer0.trafigura.uat.dlt.testnet.ch: { }
        peer0.traxys.uat.dlt.testnet.ch: { }
        peer0.umicore.uat.dlt.testnet.ch: { }
```

Все перечисленные пиры мы должны также указать в файле конфигурации config_test.yaml в
блоке peers 56 строка config_test_tmp.yaml
Укажем в файле config_test.yaml организации которые есть в нашей сети hlf. 23 строка в файле
config_test_tmp.yaml
Далее необходимо скопировать название той организации от которой будем выполнять запросы
и указать название организации в файле config_test.yaml на 5 строке.
В моем примере это организация называется **testnet**
organization: testnet
3.3.4. проверьте что все пути к файлам и дирректориям указанные в файле config_test.yaml
верные (пример папки crypto во вложении)
3.3.5. config_test_tmp.yaml - можно удалить, он больше нам не нужен т.к настройки указаны в
файле config_test.yaml

# 4. Обновление публичного ключа в hlf

Администратор отправляет запрос в hlf для изменения публичного ключа. Пользователя для hlf можно указать через параметр -u

```shell
./cli sendRequest acl changePublicKey $validatorPublicKeys -u User14
```

Проверка публичного ключа в hlf
Администратор проверяет измениться публичный ключ для адреса или нет выполнив следующие команды.

```shell
OLD_PublicKey="CTmpLBcWAtikpFYDwPkSPeQpKZALxpaGG7r5AYMiCjbG"
./cli query acl checkKeys $OLD_PublicKey -u User1

NEW_PublicKey="BREP5CVURcJ6CoTdUpJdSNgZThMXvubFueyviSHDRW4Y"
./cli query acl checkKeys $NEW_PublicKey -u User1
```