## Task done

Для запуска тестов вызывать команду:

```shell
go test -v ./tests/...
```

Для проверки отсутствия **race-conditions** добавить флаг `-race`:

```shell
go test -v ./tests/... -race -count=1
```

Для просмотра процентов __(92.2%)__ тестового покрытия вызвать:

```shell
go test -v -coverprofile cover.out ./... -count=1
go tool cover -html=cover.out
```

Для запуска бенчмарков вызвать:

```shell
 go test -v ./benchmarks/... -bench=.
# Для запуска определённого бенчмарка поставьте его название вместо точки: -bench=<name bench>
```

## Цель работы

Реализовать и сравнить 3 версии стека:

- `Simple stack`
- `Treiber stack`
- `Treiber back-off elimination stack`

Выяснить, насколько эффективна оптимизация с элиминацией по сравнению с неоптимизированной версией и как сильно она
отстаёт от простого
стека, если мы не собираемся использовать его в многопоточной программе.

## Эксперимент

### Характеристики вычислительной машины

- Процессор — AMD Ryzen 5 5500U 2.1 GHz
- Оперативная память — 16 Gb 3200 MHz
- Операционная система — Ubuntu 22.04.4 LTS

---

### Сценарии

1. Запустим 1_000_000 `Push`, а затем столько же `Pop` последовательно.
2. Запустим 10_000 `Push` на 100 горутинах, ожидая пока все закончат, а затем `Pop` с такими же параметрами.
3. Изменим 2-ой сценарий и посмотрим как изменяться цифры, если запустить 1_000_000 горутин, каждая из которых сначала
   сделает
   один `Push` и после барьера сделает `Pop`.
4. Затем рассмотрим работает ли стек Трайбера с оптимизацией быстрее на специально подобранных данных, каждая из 10_000
   горутин будет делать `Push` и `Pop` без барьеров 10_000 раз
5. Изменим сценарий 4 и посмотрим, насколько ухудшится оптимизация, если операции `Push` и `Pop` мы будем выбирать в
   рандомном порядке.

### Результаты

**Сценарий 1**

| NonConcurrent        | Time for one iteration _(ns/op)_ | Standard deviation | acceleration percentages __(%)__ | 
|----------------------|----------------------------------|--------------------|----------------------------------|
| Simple stack         | 48618067.6                       | 989054.2           | 100                              |
| Treiber stack        | 56968991.0                       | 640815.1           | 85.3                             |
| Treiber optimization | 84650766.6                       | 4953769.4          | 57.4                             | 

----
**Сценарий 2**

| LittleConcurrent     | Time for one iteration _(ns/op)_ | Standard deviation | acceleration percentages __(%)__ | 
|----------------------|----------------------------------|--------------------|----------------------------------|
| Treiber stack        | 125996729.6                      | 899816.7           | 100                              |
| Treiber optimization | 152838145.8                      | 3477201.5          | 82.4                             |

---

**Сценарий 3**

| AllConcurrent        | Time for one iteration _(ns/op)_ | Standard deviation | acceleration percentages __(%)__ | 
|----------------------|----------------------------------|--------------------|----------------------------------|
| Treiber stack        | 616976436.6                      | 10990880.4         | 100                              |
| Treiber optimization | 747104927.0                      | 42165227.9         | 82.6                             |

---

**Сценарий 4**

| PushAndPopInRow      | Time for one iteration _(ns/op)_ | Standard deviation | acceleration percentages __(%)__ | 
|----------------------|----------------------------------|--------------------|----------------------------------|
| Treiber stack        | 147294554.0                      | 3192428.3          | 100                              |
| Treiber optimization | 131419447.4                      | 1560897.5          | 112                              |

---

**Сценарий 5**

| PushAndPopRandom     | Time for one iteration _(ns/op)_ | Standard deviation | acceleration percentages __(%)__ | 
|----------------------|----------------------------------|--------------------|----------------------------------|
| Treiber stack        | 156998454.2                      | 3061606.7          | 100                              |
| Treiber optimization | 121853269.2                      | 3531401.2          | 128.8                            |

---

### Выводы

1. Для последовательных программ лучше всего подходит простой стек, однако стек Трайбера без оптимизации оказался
   медленнее всего на 15%.
2. На 10_000 или на 1_000_000 горутин, стек Трайбера с оптимизацией при сценарии, когда мы сначала конкурентно делаем
   `Push`, а только потом `Pop`, оказался медленнее стека Трайбера без оптимизации на 18%. Можно сделать вывод, что если
   логика вашей конкурентной программы подразумевает сначала общий `Push`, а только потом `Pop` - стек Трайбера с
   оптимизацией использовать не стоит.
3. Сценарии 4 и 5 очень похожи тем, что при рандомном `Push` или `Pop` стек Трайбера с оптимизацией оказался быстрее 
   своего предка без оптимизации на 28%. Думаю, это число может быть ещё больше, на более мощной вычислительной машине.

## Материалы

- The Art of Multiprocessor Programming (Chapter 11)
- A Scalable Lock-free Stack Algorithm 2004 (Danny Hendler, Nir Shavit, Lena Yerushalmi)
- Неплохая [статья](https://max-inden.de/post/2020-03-28-elimination-backoff-stack/) о lock-free stack
