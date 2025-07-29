package service

import (
	"bytes"
	"fmt"
	"github.com/nescool101/rentManager/storage"
	"html/template"
	"log"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"

	"github.com/nescool101/rentManager/model"
)

var payers []model.Payer

// Load all payers from storage
func LoadPayers() {
	p, err := storage.GetPayers()
	if err != nil {
		log.Fatal("❌ [ERROR] Failed to load payers:", err)
	}
	payers = p
	log.Printf("✅ [INFO] Successfully loaded %d payers from storage.", len(payers))
}

// Notify all payers about rent reminders
func NotifyAll() {
	loc, _ := time.LoadLocation("America/New_York") // Load EST timezone
	today := time.Now().In(loc)                     // Convert to EST

	log.Printf("✅ [INFO] loaded %d payers from storage.", len(payers))
	for _, payer := range payers {
		if payer.RentalEmail == "" {
			continue
		}

		rentalDay := payer.RentalDate.Day()
		rentalMonth := payer.RentalDate.Month()
		rentalYear := payer.RentalDate.Year()

		// One-month reminder
		sameMonthReminder(today, rentalDay, payer)

		// One-year anniversary reminder
		sameYearReminder(today, rentalDay, rentalMonth, rentalYear, payer)
	}
}

// Send one-year rental anniversary reminder
func sameYearReminder(today time.Time, rentalDay int, rentalMonth time.Month, rentalYear int, payer model.Payer) {
	if today.Day() == rentalDay && today.Month() == rentalMonth && today.Year() != rentalYear {
		log.Printf("📩 [1-YEAR ANNIVERSARY] Sending to Tenant: %s (%s) and Renter: %s (%s)",
			payer.Name, payer.RentalEmail, payer.RenterName, payer.RenterEmail)

		subject := "🏡 Aniversario de Arrendamiento"
		body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
		    <meta charset="UTF-8">
		    <title>Aniversario de Arrendamiento</title>
		    <style>
		        body { font-family: Arial, sans-serif; }
		        .container { padding: 20px; }
		        .highlight { font-weight: bold; color: #007BFF; }
		    </style>
		</head>
		<body>
		    <div class="container">
		        <h2>🏡 ¡Feliz Aniversario de Arrendamiento, %s!</h2>
		        <p>Hoy se cumple un año desde que inició su contrato de arrendamiento para la propiedad en:</p>
		        <p class="highlight">%s</p>
		        <p>Le agradecemos su confianza y esperamos que su experiencia haya sido excelente.</p>
		        <p>¿Desea renovar su contrato de arrendamiento?</p>
		        <p>Por favor, comuníquese con nosotros para discutir las opciones de renovación.</p>
		        <hr>
		        <p>Atentamente,</p>
		        <p><strong>%s</strong></p>
		    </div>
		</body>
		</html>
		`, payer.Name, payer.PropertyAddress, payer.RenterName)

		// Send email to Tenant
		errTenant := sendSimpleEmail(payer.RentalEmail, subject, body)
		if errTenant != nil {
			log.Printf("❌ [FAILED] 1-Year Anniversary Email NOT sent to Tenant: %s (%s) - Error: %v",
				payer.Name, payer.RentalEmail, errTenant)
		} else {
			log.Printf("✅ [SENT] 1-Year Anniversary Email sent to Tenant: %s (%s)",
				payer.Name, payer.RentalEmail)
		}

		// Send email to Renter (if available)
		if payer.RenterEmail != "" {
			errRenter := sendSimpleEmail(payer.RenterEmail, subject, body)
			if errRenter != nil {
				log.Printf("❌ [FAILED] 1-Year Anniversary Email NOT sent to Renter: %s (%s) - Error: %v",
					payer.RenterName, payer.RenterEmail, errRenter)
			} else {
				log.Printf("✅ [SENT] 1-Year Anniversary Email sent to Renter: %s (%s)",
					payer.RenterName, payer.RenterEmail)
			}
		}
	}
}

// Send one-month rental reminder
func sameMonthReminder(today time.Time, rentalDay int, payer model.Payer) {
	if today.Day() == rentalDay {
		log.Printf("📩 [1-MONTH REMINDER] Sending to: %s (%s)", payer.Name, payer.RentalEmail)
		err := sendEmail(payer.RentalEmail, payer)
		if err != nil {
			log.Printf("❌ [FAILED] 1-Month Reminder NOT sent to %s (%s) - Error: %v", payer.Name, payer.RentalEmail, err)
			return
		}
		log.Printf("✅ [SENT] 1-Month Reminder sent to: %s (%s)", payer.Name, payer.RentalEmail)
	}
	log.Printf("❌ not [SENT] 1-Month Reminder sent to: %s (%s) (%s)", payer.Name, payer.RentalEmail, rentalDay)
}

// EmailTemplate represents the structure of the email data
type EmailTemplate struct {
	EmisorNombre         string
	EmisorNIT            string
	EmisorDireccion      string
	EmisorTelefono       string
	EmisorEmail          string
	NumeroCuenta         int
	FechaEmision         string
	ArrendatarioNombre   string
	ArrendatarioNIT      string
	InmuebleDireccion    string
	TipoInmueble         string
	FechaInicio          string
	FechaFinal           string
	ValorMensual         string
	Subtotal             string
	TotalPagar           string
	CondicionesPago      string
	Banco                string
	TipoCuenta           string
	NumeroCuentaBancaria string
	TitularCuenta        string
	Observaciones        string
	ArrendadorNombre     string
	UnpaidMonths         int
	TotalDue             string
}

// Email template in HTML format
const emailTemplateHTML = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Cuenta de Cobro</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f4f4f4;
        }
        .container {
            width: 100%;
            max-width: 600px;
            margin: 20px auto;
            background-color: #ffffff;
            border-radius: 8px;
            box-shadow: 0 4px 8px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .header {
            background-color: #4CAF50;
            color: #ffffff;
            padding: 20px;
            text-align: center;
        }
        .header h2 {
            margin: 0;
        }
        .content {
            padding: 20px;
        }
        .section {
            margin-bottom: 20px;
        }
        .section h4 {
            border-bottom: 2px solid #eeeeee;
            padding-bottom: 5px;
            margin-bottom: 10px;
            color: #333333;
        }
        .info-grid {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 10px;
        }
        .info-grid p {
            margin: 5px 0;
        }
        .table-container {
            margin-top: 20px;
        }
        table {
            width: 100%;
            border-collapse: collapse;
        }
        th, td {
            border: 1px solid #dddddd;
            text-align: left;
            padding: 8px;
        }
        th {
            background-color: #f2f2f2;
        }
        .total {
            text-align: right;
            font-size: 1.2em;
            font-weight: bold;
            margin-top: 20px;
        }
        .warning {
            background-color: #ffebee;
            border: 1px solid #ffcdd2;
            padding: 15px;
            border-radius: 4px;
            margin-top: 20px;
        }
        .warning h3 {
            color: #c62828;
            margin-top: 0;
        }
        .footer {
            background-color: #f8f8f8;
            padding: 20px;
            text-align: center;
            font-size: 0.9em;
            color: #555555;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>Cuenta de Cobro de Arrendamiento</h2>
        </div>
        <div class="content">
            <div class="section">
                <h4>Información del Arrendador</h4>
                <div class="info-grid">
                    <p><strong>Nombre:</strong> {{.EmisorNombre}}</p>
                    <p><strong>NIT/Cédula:</strong> {{.EmisorNIT}}</p>
                    <p><strong>Dirección:</strong> {{.EmisorDireccion}}</p>
                    <p><strong>Teléfono:</strong> {{.EmisorTelefono}}</p>
                </div>
            </div>

            <div class="section">
                <h4>Información del Arrendatario</h4>
                <div class="info-grid">
                    <p><strong>Nombre:</strong> {{.ArrendatarioNombre}}</p>
                    <p><strong>NIT/Cédula:</strong> {{.ArrendatarioNIT}}</p>
                </div>
                <p><strong>Dirección del Inmueble:</strong> {{.InmuebleDireccion}}</p>
            </div>

            <div class="section">
                <h4>Detalles de la Cuenta</h4>
                <p><strong>Número de Cuenta de Cobro:</strong> {{.NumeroCuenta}}</p>
                <p><strong>Fecha de Emisión:</strong> {{.FechaEmision}}</p>
            </div>

            <div class="table-container">
                <table>
                    <tr>
                        <th>Descripción</th>
                        <th>Periodo</th>
                        <th>Valor</th>
                    </tr>
                    <tr>
                        <td>Canon de Arrendamiento - {{.TipoInmueble}}</td>
                        <td>{{.FechaInicio}} a {{.FechaFinal}}</td>
                        <td>{{.ValorMensual}}</td>
                    </tr>
                </table>
            </div>

            <div class="total">
                <p>Subtotal: {{.Subtotal}}</p>
                <p>Total a Pagar: {{.TotalPagar}}</p>
            </div>

            {{if gt .UnpaidMonths 0}}
            <div class="warning">
                <h3>⚠️ Recordatorio de Pago Atrasado</h3>
                <p>Hemos notado que tienes <strong>{{.UnpaidMonths}} mes(es)</strong> de arriendo pendientes.</p>
                <p>El monto total adeudado es de <strong>{{.TotalDue}}</strong>.</p>
                <p>Por favor, realiza el pago a la brevedad posible para evitar inconvenientes.</p>
            </div>
            {{end}}

            <div class="section">
                <h4>Instrucciones de Pago</h4>
                <p><strong>Condiciones:</strong> {{.CondicionesPago}}</p>
                <p><strong>Banco:</strong> {{.Banco}}</p>
                <p><strong>Tipo de Cuenta:</strong> {{.TipoCuenta}}</p>
                <p><strong>Número de Cuenta:</strong> {{.NumeroCuentaBancaria}}</p>
                <p><strong>Titular:</strong> {{.TitularCuenta}}</p>
            </div>

            <div class="section">
                <h4>Observaciones Adicionales</h4>
                <p>{{.Observaciones}}</p>
            </div>
        </div>
        <div class="footer">
            <p>Este es un recordatorio de pago automático. Si ya has realizado el pago, por favor ignora este mensaje.</p>
            <p>Atentamente,<br><strong>{{.ArrendadorNombre}}</strong></p>
            <p>{{.EmisorEmail}}</p>
        </div>
    </div>
</body>
</html>
`

func sendSimpleEmail(to, subject, body string) error {
	host := "smtp.gmail.com"
	portStr := "587"
	user := "nescool10001@gmail.com"
	pass := "bndpfcmeoyhhudyz"

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("❌ [ERROR] Invalid SMTP port: %s", portStr)
		return err
	}

	d := gomail.NewDialer(host, port, user, pass)
	m := gomail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body) // Send HTML email

	err = d.DialAndSend(m)
	if err != nil {
		log.Printf("❌ [EMAIL NOT SENT] %s (%s) - Error: %v", to, err)
		return err
	}

	log.Printf("✅ [EMAIL SENT] %s", to)
	return nil
}

func sendEmail(to string, payer model.Payer) error {
	//host := os.Getenv("MAIL_HOST")
	//portStr := os.Getenv("MAIL_PORT")
	//user := os.Getenv("MAIL_USERNAME")
	//pass := os.Getenv("MAIL_PASSWORD")
	host := "smtp.gmail.com"
	portStr := "587"
	user := "nescool10001@gmail.com"
	pass := "bndpfcmeoyhhudyz"

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}

	// Convert MonthlyRent to an integer (removing "USD" or currency text)
	totalDue := 0
	if payer.UnpaidMonths > 0 {
		totalDue = payer.MonthlyRent * payer.UnpaidMonths
	}

	data := EmailTemplate{
		EmisorNombre:         "Mi Empresa S.A.",
		EmisorNIT:            "123456789",
		EmisorDireccion:      "Calle 123, Ciudad",
		EmisorTelefono:       "555-1234",
		EmisorEmail:          "empresa@example.com",
		NumeroCuenta:         rentalDateToInt(payer.RentalDate),
		FechaEmision:         payer.RentalDate.Format("02/01/2006"),
		ArrendatarioNombre:   payer.Name,
		ArrendatarioNIT:      payer.NIT,
		InmuebleDireccion:    payer.PropertyAddress,
		TipoInmueble:         payer.PropertyType,
		FechaInicio:          payer.RentalStart.Format("02/01/2006"),
		FechaFinal:           payer.RentalEnd.Format("02/01/2006"),
		ValorMensual:         strconv.Itoa(payer.MonthlyRent),
		Subtotal:             strconv.Itoa(payer.MonthlyRent),
		TotalPagar:           strconv.Itoa(payer.MonthlyRent),
		CondicionesPago:      "Pago antes del 5 de cada mes",
		Banco:                payer.BankName,
		TipoCuenta:           payer.AccountType,
		NumeroCuentaBancaria: payer.BankAccountNumber,
		TitularCuenta:        payer.AccountHolder,
		Observaciones:        payer.AdditionalNotes,
		ArrendadorNombre:     payer.RenterName,
		UnpaidMonths:         payer.UnpaidMonths,              // Use directly
		TotalDue:             strconv.Itoa(totalDue) + " COP", // Only if unpaid months exist
	}

	// Parse and execute the HTML template
	tmpl, err := template.New("email").Parse(emailTemplateHTML)
	if err != nil {
		log.Println("Error parsing template:", err)
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		log.Println("Error executing template:", err)
		return err
	}

	d := gomail.NewDialer(host, port, user, pass)
	m := gomail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Cuenta de Cobro Arrendamiento")
	m.SetBody("text/html", body.String()) // Send HTML email

	log.Printf("✅ [EMAIL SENT] %s (%s)", payer.Name, to)
	return d.DialAndSend(m)
}

func rentalDateToInt(date time.Time) int {
	return date.Year()*10000 + int(date.Month())*100 + date.Day()
}
