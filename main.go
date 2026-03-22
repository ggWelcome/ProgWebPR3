package main

import (
	"html/template" // Пакет для роботи з HTML-шаблонами
	"log"           // Пакет для логування повідомлень
	"net/http"      // Пакет для створення веб-сервера та обробки HTTP-запитів
	"strconv"       // Пакет для конвертації рядків у числа
	"time"          // Пакет для роботи з датами та часом
)

// Структура для збереження даних про споживача
type Consumer struct {
	Name      string    // Ім’я споживача
	Address   string    // Адреса споживача
	Connected time.Time // Дата підключення
}

// Структура для збереження даних датчика
type SensorData struct {
	Voltage float64 // Напруга у вольтах
	Current float64 // Струм у амперах
	Power   float64 // Потужність у ватах (обчислюється як Voltage * Current)
}

// Глобальні змінні для збереження списку споживачів та даних датчиків
var consumers []Consumer
var sensors []SensorData

// Обробник головної сторінки з формою
func formHandler(w http.ResponseWriter, r *http.Request) {
	// Завантаження HTML-шаблону form.html
	tmpl := template.Must(template.ParseFiles("form.html"))
	// Виведення шаблону у відповідь користувачу
	tmpl.Execute(w, nil)
}

// Обробник форми після натискання кнопки "Зберегти"
func submitHandler(w http.ResponseWriter, r *http.Request) {
	// Перевірка, чи запит є POST (щоб не обробляти GET-запити)
	if r.Method != http.MethodPost {
		// Якщо метод не POST, перенаправляємо користувача назад на головну сторінку
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Отримуємо значення прихованого поля formType (визначає тип форми: consumer або sensor)
	formType := r.FormValue("formType")

	// Виконуємо різні дії залежно від типу форми
	switch formType {
	case "consumer": // Якщо користувач додає споживача
		name := r.FormValue("name")       // Отримуємо ім’я
		address := r.FormValue("address") // Отримуємо адресу
		dateStr := r.FormValue("date")    // Отримуємо дату підключення у вигляді рядка

		// Конвертуємо рядок у формат дати (YYYY-MM-DD)
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			// Якщо формат дати неправильний — повертаємо помилку
			http.Error(w, "Невірний формат дати!", http.StatusBadRequest)
			return
		}

		// Створюємо новий об’єкт споживача
		consumer := Consumer{name, address, date}
		// Додаємо його у глобальний список споживачів
		consumers = append(consumers, consumer)

		// Завантажуємо шаблон customer.html
		tmpl := template.Must(template.ParseFiles("customer.html"))
		// Виводимо сторінку з даними про споживача
		tmpl.Execute(w, consumer)

	case "sensor": // Якщо користувач додає дані датчика
		// Отримуємо значення напруги та струму з форми і конвертуємо у числа
		voltage, err1 := strconv.ParseFloat(r.FormValue("voltage"), 64)
		current, err2 := strconv.ParseFloat(r.FormValue("current"), 64)

		// Перевірка на помилки або некоректні значення (напруга та струм не можуть бути від’ємними)
		if err1 != nil || err2 != nil || voltage < 0 || current < 0 {
			http.Error(w, "Некоректні дані датчика!", http.StatusBadRequest)
			return
		}

		// Обчислюємо потужність (Вт = В * А)
		power := voltage * current
		// Створюємо новий об’єкт SensorData
		sensor := SensorData{voltage, current, power}
		// Додаємо його у глобальний список датчиків
		sensors = append(sensors, sensor)

		// Завантажуємо шаблон electricity.html
		tmpl := template.Must(template.ParseFiles("electricity.html"))
		// Виводимо сторінку з даними датчика
		tmpl.Execute(w, sensor)
	}
}

// Головна функція програми
func main() {
	// Реєструємо обробник для головної сторінки (форма)
	http.HandleFunc("/", formHandler)
	// Реєструємо обробник для відправки даних форми
	http.HandleFunc("/submit", submitHandler)

	// Виводимо повідомлення у консоль про запуск сервера
	log.Println("Сервер запущено на http://localhost:8080")
	// Запускаємо веб-сервер на порту 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
