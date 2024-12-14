package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

type QRData struct {
	UniqueID   string `json:"uniqueId"`
	BatchID    int    `json:"batchId"`
	DeviceName string `json:"deviceName"`
}

func main() {
	file, err := os.Open("qrcodes.txt")
	if err != nil {
		fmt.Printf("Error opening text file: %v\n", err)
		return
	}
	defer file.Close()

	var qrDataStrings []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var data QRData
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			fmt.Printf("Error parsing line: %v\n", err)
			continue
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("Error creating JSON: %v\n", err)
			continue
		}

		qrDataStrings = append(qrDataStrings, string(jsonData))
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	generateQRCodes(qrDataStrings)
}

func generateQRCodes(dataStrings []string) {
	qrDir := "generated-qrs"
	if err := os.MkdirAll(qrDir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	const (
		qrSize      = 50
		marginX     = 15
		marginY     = 10
		spacing     = 10
		codesPerRow = 3
	)

	currentX := marginX
	currentY := marginY
	count := 0

	for _, dataString := range dataStrings {
		var data QRData
		if err := json.Unmarshal([]byte(dataString), &data); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			continue
		}

		tempFile := filepath.Join(qrDir, fmt.Sprintf("%s_qr.png", data.UniqueID))
		qrc, err := qrcode.NewWith(
			dataString,
			qrcode.WithEncodingMode(qrcode.EncModeByte),
			qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionQuart),
		)
		if err != nil {
			fmt.Printf("Could not generate QRCode: %v\n", err)
			continue
		}

		w, err := standard.New(tempFile,
			standard.WithQRWidth(uint8(qrSize)),
			standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
		)
		if err != nil {
			fmt.Printf("standard.New failed: %v\n", err)
			continue
		}

		if err = qrc.Save(w); err != nil {
			fmt.Printf("Could not save image: %v\n", err)
			continue
		}

		pdf.Image(tempFile, float64(currentX), float64(currentY), qrSize, qrSize, false, "", 0, "")

		pdf.SetFont("Arial", "", 8)
		textWidth := pdf.GetStringWidth(data.UniqueID)
		textX := float64(currentX) + (qrSize-textWidth)/2
		pdf.Text(textX, float64(currentY)+qrSize+5, data.UniqueID)

		count++
		currentX += qrSize + spacing

		if count%codesPerRow == 0 {
			currentX = marginX
			currentY += qrSize + spacing + 10

			if currentY > 250 {
				pdf.AddPage()
				currentY = marginY
			}
		}

		os.Remove(tempFile)
	}

	pdfPath := filepath.Join(qrDir, "qr_codes.pdf")
	if err := pdf.OutputFileAndClose(pdfPath); err != nil {
		fmt.Printf("Error saving PDF: %v\n", err)
		return
	}

	fmt.Printf("Successfully generated PDF with QR codes: %s\n", pdfPath)
}
