# go-grepdirrec

Программа на Go для рекурсивного поиска указанного слова в файлах и директориях с ограничением числа параллельных горутин.

## Использование

```bash
go run grepdirrec.go <слово> <путь_к_файлу_или_директории>
```

## Пример
```bash
go run grepdirrec.go bubble ./example.txt
```
## Результат
✅ WORD "bubble" was found in FILE: ./Projects/go/example.txt  
✅ WORD "bubble" was found in FILE: ./Projects/go/subdir/test.txt
