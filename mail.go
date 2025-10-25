package main

import (
    "crypto/tls"
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/wneessen/go-mail"
)

// Структуры данных
type MailConfig struct {
    Host        string
    Port        int
    Username    string
    Password    string
    To          string
    ToFile      string
    Subject     string
    Body        string
    BodyFile    string
    IsHTML      bool
    Attachments string
    MessageID   string
}

type EmailLog struct {
    ID         string            `json:"id"`
    Timestamp  time.Time         `json:"timestamp"`
    From       string            `json:"from"`
    Subject    string            `json:"subject"`
    Status     string            `json:"overall_status"`
    Recipients []RecipientStatus `json:"recipients"`
}

type RecipientStatus struct {
    Email   string    `json:"email"`
    Status  string    `json:"status"`
    Time    time.Time `json:"time"`
    Error   string    `json:"error,omitempty"`
    Message string    `json:"message,omitempty"`
}

type RecipientsFile struct {
    Emails []string `json:"emails"`
}

type SendEmailRequest struct {
    To          []string `json:"to" binding:"required"`
    Subject     string   `json:"subject" binding:"required"`
    Body        string   `json:"body"`
    BodyFile    string   `json:"body_file,omitempty"`
    IsHTML      bool     `json:"is_html,omitempty"`
    Attachments []string `json:"attachments,omitempty"`
    MessageID   string   `json:"message_id,omitempty"`
}

type SendEmailResponse struct {
    ID      string `json:"id"`
    Status  string `json:"status"`
    Message string `json:"message"`
}

type StatusResponse struct {
    ID        string            `json:"id"`
    Timestamp time.Time         `json:"timestamp"`
    Status    string            `json:"overall_status"`
    Details   []RecipientStatus `json:"details,omitempty"`
}

// Глобальные переменные
var (
    smtpHost     string
    smtpPort     int
    smtpUsername string
    smtpPassword string
    apiPort      string
    
    // CLI флаги
    cliTo          string
    cliToFile      string
    cliSubject     string
    cliBody        string
    cliBodyFile    string
    cliIsHTML      bool
    cliAttachments string
    cliMessageID   string
)

func main() {
    // Парсим флаги
    parseFlags()

    // Если указан порт API, запускаем сервер
    if apiPort != "" {
        startAPIServer()
    } else {
        // Иначе работаем в CLI режиме
        runCLI()
    }
}

func parseFlags() {
    // Основные флаги
    flag.StringVar(&smtpHost, "host", "mail.example.com", "SMTP хост")
    flag.IntVar(&smtpPort, "port", 25, "SMTP порт")
    flag.StringVar(&smtpUsername, "user", "archive@mail.example.com", "Имя пользователя")
    flag.StringVar(&smtpPassword, "pass", "StrongPass", "Пароль")
    flag.StringVar(&apiPort, "api", "", "REST API порт (например :8080)")
    
    // CLI флаги
    flag.StringVar(&cliTo, "to", "", "Получатель (обязательно, несколько через запятую)")
    flag.StringVar(&cliToFile, "toFile", "", "JSON файл со списком получателей")
    flag.StringVar(&cliSubject, "subject", "Тестовое письмо", "Тема письма")
    flag.StringVar(&cliBody, "body", "", "Текст письма")
    flag.StringVar(&cliBodyFile, "bodyFile", "", "Файл с содержимым письма")
    flag.BoolVar(&cliIsHTML, "html", false, "HTML формат")
    flag.StringVar(&cliAttachments, "attach", "", "Файлы для вложения (через запятую)")
    flag.StringVar(&cliMessageID, "id", "", "ID письма (автогенерация если не указан)")

    flag.Parse()
}

func runCLI() {
    config := &MailConfig{
        Host:        smtpHost,
        Port:        smtpPort,
        Username:    smtpUsername,
        Password:    smtpPassword,
        To:          cliTo,
        ToFile:      cliToFile,
        Subject:     cliSubject,
        Body:        cliBody,
        BodyFile:    cliBodyFile,
        IsHTML:      cliIsHTML,
        Attachments: cliAttachments,
        MessageID:   cliMessageID,
    }

    // Проверяем обязательные поля
    if config.To == "" && config.ToFile == "" {
        fmt.Println("Ошибка: необходимо указать получателей через -to или -toFile")
        fmt.Println("Использование в CLI режиме:")
        fmt.Println("  ./go-mail -to email@example.com -subject 'Тема' -body 'Текст'")
        fmt.Println("  ./go-mail -toFile recipients.json -subject 'Тема' -body 'Текст'")
        fmt.Println("\nВсе параметры:")
        flag.PrintDefaults()
        os.Exit(1)
    }

    // Создаем папку для логов если не существует
    if err := os.MkdirAll("logs", 0755); err != nil {
        log.Fatal("❌ Ошибка создания папки logs: ", err)
    }

    // Генерируем ID письма если не передан
    if config.MessageID == "" {
        config.MessageID = generateMessageID()
    }

    fmt.Printf("📧 Отправка письма ID: %s\n", config.MessageID)
    fmt.Printf("   Сервер: %s:%d\n", config.Host, config.Port)
    fmt.Printf("   От: %s\n", config.Username)
    fmt.Printf("   Тема: %s\n", config.Subject)

    // Получаем список получателей
    recipients, err := getRecipients(config)
    if err != nil {
        log.Fatal("❌ Ошибка получения списка получателей: ", err)
    }

    if len(recipients) == 0 {
        log.Fatal("❌ Не указаны получатели")
    }

    fmt.Printf("   Получатели: %s\n", strings.Join(recipients, ", "))

    // Создаем лог-файл
    logEntry := &EmailLog{
        ID:         config.MessageID,
        Timestamp:  time.Now(),
        From:       config.Username,
        Subject:    config.Subject,
        Recipients: make([]RecipientStatus, 0, len(recipients)),
    }

    // Если один получатель - синхронная отправка
    // Если несколько - асинхронная
    if len(recipients) == 1 {
        if err := sendEmail(config, recipients[0], logEntry); err != nil {
            log.Fatal("❌ Ошибка отправки: ", err)
        }
    } else {
        sendEmailsAsync(config, recipients, logEntry)
    }

    // Определяем общий статус
    if allSuccess(logEntry.Recipients) {
        logEntry.Status = "success"
    } else if anySuccess(logEntry.Recipients) {
        logEntry.Status = "partial_success"
    } else {
        logEntry.Status = "failed"
    }

    // Сохраняем финальный лог
    if err := saveLog(logEntry); err != nil {
        log.Printf("❌ Ошибка сохранения лога: %v", err)
    }

    fmt.Printf("✅ Процесс отправки завершен! Статус: %s\n", logEntry.Status)
    fmt.Printf("📁 Лог сохранен: logs/%s.json\n", logEntry.ID)
}

func startAPIServer() {
    // Создаем папку для логов если не существует
    if err := os.MkdirAll("logs", 0755); err != nil {
        log.Fatal("❌ Ошибка создания папки logs: ", err)
    }

    // Настраиваем Gin
    gin.SetMode(gin.ReleaseMode)
    router := gin.Default()

    // Middleware для логирования
    router.Use(gin.Logger())

    // Роуты API
    router.POST("/api/send", handleSendEmail)
    router.GET("/api/status/:id", handleGetStatus)
    router.GET("/api/health", handleHealthCheck)

    fmt.Printf("🚀 REST API сервер запущен на порту %s\n", apiPort)
    fmt.Printf("📧 Эндпоинты:\n")
    fmt.Printf("   POST /api/send     - Отправка email\n")
    fmt.Printf("   GET  /api/status/:id - Статус отправки\n")
    fmt.Printf("   GET  /api/health   - Проверка здоровья\n")

    if err := router.Run(apiPort); err != nil {
        log.Fatal("❌ Ошибка запуска API сервера: ", err)
    }
}

func handleSendEmail(c *gin.Context) {
    var req SendEmailRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат запроса: " + err.Error(),
        })
        return
    }

    // Генерируем ID если не указан
    if req.MessageID == "" {
        req.MessageID = generateMessageID()
    }

    // Создаем конфиг для отправки
    config := &MailConfig{
        Host:        smtpHost,
        Port:        smtpPort,
        Username:    smtpUsername,
        Password:    smtpPassword,
        Subject:     req.Subject,
        Body:        req.Body,
        BodyFile:    req.BodyFile,
        IsHTML:      req.IsHTML,
        Attachments: strings.Join(req.Attachments, ","),
        MessageID:   req.MessageID,
    }

    // Создаем лог
    logEntry := &EmailLog{
        ID:         req.MessageID,
        Timestamp:  time.Now(),
        From:       smtpUsername,
        Subject:    req.Subject,
        Recipients: make([]RecipientStatus, 0, len(req.To)),
    }

    // Запускаем асинхронную отправку
    go sendEmailsAsync(config, req.To, logEntry)

    // Возвращаем немедленный ответ
    c.JSON(http.StatusAccepted, SendEmailResponse{
        ID:      req.MessageID,
        Status:  "accepted",
        Message: "Письмо принято в обработку",
    })
}

func handleGetStatus(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "ID обязателен",
        })
        return
    }

    // Пытаемся прочитать лог файл
    logEntry, err := readLog(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Лог не найден",
        })
        return
    }

    // Создаем ответ
    response := StatusResponse{
        ID:        logEntry.ID,
        Timestamp: logEntry.Timestamp,
        Status:    logEntry.Status,
        Details:   logEntry.Recipients,
    }

    c.JSON(http.StatusOK, response)
}

func handleHealthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":    "healthy",
        "timestamp": time.Now(),
        "service":   "Email Sender API",
    })
}

func readLog(id string) (*EmailLog, error) {
    filename := fmt.Sprintf("%s.json", id)
    filepath := filepath.Join("logs", filename)

    data, err := os.ReadFile(filepath)
    if err != nil {
        return nil, err
    }

    var logEntry EmailLog
    if err := json.Unmarshal(data, &logEntry); err != nil {
        return nil, err
    }

    return &logEntry, nil
}

func getRecipients(config *MailConfig) ([]string, error) {
    var recipients []string

    // Добавляем получателей из флага -to
    if config.To != "" {
        toRecipients := parseRecipients(config.To)
        recipients = append(recipients, toRecipients...)
    }

    // Добавляем получателей из файла
    if config.ToFile != "" {
        fileRecipients, err := loadRecipientsFromFile(config.ToFile)
        if err != nil {
            return nil, err
        }
        recipients = append(recipients, fileRecipients...)
    }

    // Удаляем дубликаты
    recipients = removeDuplicates(recipients)

    return recipients, nil
}

func loadRecipientsFromFile(filename string) ([]string, error) {
    // Читаем файл
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, fmt.Errorf("ошибка чтения файла %s: %v", filename, err)
    }

    // Парсим JSON
    var recipientsFile RecipientsFile
    if err := json.Unmarshal(data, &recipientsFile); err != nil {
        return nil, fmt.Errorf("ошибка парсинга JSON файла %s: %v", filename, err)
    }

    // Очищаем email адреса
    var cleanRecipients []string
    for _, email := range recipientsFile.Emails {
        email = strings.TrimSpace(email)
        if email != "" && isValidEmail(email) {
            cleanRecipients = append(cleanRecipients, email)
        }
    }

    if len(cleanRecipients) == 0 {
        return nil, fmt.Errorf("файл %s не содержит валидных email адресов", filename)
    }

    return cleanRecipients, nil
}

func parseRecipients(recipientsStr string) []string {
    recipients := strings.Split(recipientsStr, ",")
    var cleanRecipients []string
    for _, r := range recipients {
        r = strings.TrimSpace(r)
        if r != "" && isValidEmail(r) {
            cleanRecipients = append(cleanRecipients, r)
        }
    }
    return cleanRecipients
}

func removeDuplicates(emails []string) []string {
    seen := make(map[string]bool)
    var result []string
    for _, email := range emails {
        if !seen[email] {
            seen[email] = true
            result = append(result, email)
        }
    }
    return result
}

func isValidEmail(email string) bool {
    // Простая проверка email
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func generateMessageID() string {
    return uuid.New().String()
}

func sendEmailsAsync(config *MailConfig, recipients []string, logEntry *EmailLog) {
    var wg sync.WaitGroup
    var mu sync.Mutex

    fmt.Printf("🚀 Асинхронная отправка %d писем...\n", len(recipients))

    for _, recipient := range recipients {
        wg.Add(1)
        go func(to string) {
            defer wg.Done()
            
            recipientStatus := RecipientStatus{
                Email: to,
                Time:  time.Now(),
            }

            if err := sendEmail(config, to, logEntry); err != nil {
                recipientStatus.Status = "error"
                recipientStatus.Error = err.Error()
            } else {
                recipientStatus.Status = "success"
                recipientStatus.Message = "Письмо успешно отправлено"
            }

            // Безопасно добавляем результат в лог
            mu.Lock()
            logEntry.Recipients = append(logEntry.Recipients, recipientStatus)
            mu.Unlock()

            // Выводим прогресс в консоль
            if recipientStatus.Status == "success" {
                fmt.Printf("   ✅ Успешно: %s\n", to)
            } else {
                fmt.Printf("   ❌ Ошибка: %s - %s\n", to, recipientStatus.Error)
            }
        }(recipient)
    }

    wg.Wait()

    // Определяем общий статус после завершения всех отправок
    if allSuccess(logEntry.Recipients) {
        logEntry.Status = "success"
    } else if anySuccess(logEntry.Recipients) {
        logEntry.Status = "partial_success"
    } else {
        logEntry.Status = "failed"
    }

    // Сохраняем финальный лог
    if err := saveLog(logEntry); err != nil {
        fmt.Printf("❌ Ошибка сохранения лога: %v\n", err)
    } else {
        fmt.Printf("📁 Лог сохранен: logs/%s.json\n", logEntry.ID)
    }
}

func sendEmail(config *MailConfig, recipient string, logEntry *EmailLog) error {
    // Создаем клиент
    client, err := mail.NewClient(config.Host, 
        mail.WithPort(config.Port),
        mail.WithSMTPAuth(mail.SMTPAuthPlain),
        mail.WithUsername(config.Username),
        mail.WithPassword(config.Password),
        mail.WithTLSPolicy(mail.TLSMandatory),
    )
    if err != nil {
        return fmt.Errorf("ошибка создания клиента: %v", err)
    }

    // Настраиваем TLS конфигурацию
    tlsConfig := &tls.Config{
        InsecureSkipVerify: true,
        ServerName:         config.Host,
    }
    client.SetTLSConfig(tlsConfig)

    // Создаем сообщение
    msg := mail.NewMsg()
    
    if err := msg.From(config.Username); err != nil {
        return fmt.Errorf("ошибка установки отправителя: %v", err)
    }
    if err := msg.To(recipient); err != nil {
        return fmt.Errorf("ошибка установки получателя: %v", err)
    }
    
    msg.Subject(config.Subject)
    
    // Тело сообщения
    body := config.Body
    if config.BodyFile != "" {
        content, err := os.ReadFile(config.BodyFile)
        if err != nil {
            return fmt.Errorf("ошибка чтения файла: %v", err)
        }
        body = string(content)
    }
    
    if config.IsHTML {
        msg.SetBodyString(mail.TypeTextHTML, body)
    } else {
        msg.SetBodyString(mail.TypeTextPlain, body)
    }
    
    // Добавляем вложения если есть
    if config.Attachments != "" {
        files := strings.Split(config.Attachments, ",")
        for _, file := range files {
            file = strings.TrimSpace(file)
            if file != "" {
                msg.AttachFile(file)
            }
        }
    }

    // Отправляем
    if err := client.DialAndSend(msg); err != nil {
        return fmt.Errorf("ошибка отправки: %v", err)
    }

    return nil
}

func saveLog(logEntry *EmailLog) error {
    // Создаем имя файла на основе ID
    filename := fmt.Sprintf("%s.json", logEntry.ID)
    filepath := filepath.Join("logs", filename)

    // Форматируем JSON с отступами
    jsonData, err := json.MarshalIndent(logEntry, "", "  ")
    if err != nil {
        return fmt.Errorf("ошибка кодирования JSON: %v", err)
    }

    // Записываем в файл
    if err := os.WriteFile(filepath, jsonData, 0644); err != nil {
        return fmt.Errorf("ошибка записи файла: %v", err)
    }

    return nil
}

func allSuccess(recipients []RecipientStatus) bool {
    for _, r := range recipients {
        if r.Status != "success" {
            return false
        }
    }
    return true
}

func anySuccess(recipients []RecipientStatus) bool {
    for _, r := range recipients {
        if r.Status == "success" {
            return true
        }
    }
    return false
}
