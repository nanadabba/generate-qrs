# QR Code Generator

A Go program that generates QR codes from JSON data and arranges them in a PDF document.

## Description

This program reads JSON data from a text file, generates QR codes for each entry, and creates a PDF document with the QR codes arranged in a 3-column layout. Each QR code includes its unique identifier printed below it.

## Prerequisites

-   Go 1.16 or higher
-   The following Go packages:
    ```bash
    go get github.com/jung-kurt/gofpdf
    go get github.com/yeqown/go-qrcode/v2
    ```

## Installation

1. Clone this repository:

    ```bash
    git clone <repository-url>
    cd <repository-name>
    ```

2. Install dependencies:
    ```bash
    go mod tidy
    ```

## Usage

1. Create a `qrcodes.txt` file in the project root directory with your JSON data. Each line should contain a JSON object with the following format:

    ```json
    { "unique_id": "your-unique-id", "batch_id": number }
    ```

    A sample file (`qrcodes_sample.txt`) is provided for reference, showing the expected data format:

    ```json
    {
        "unique_id": "dabbadabbadabbadabbadabbadabbadabbadabba",
        "batch_id": 41221
    }
    ```

2. Run the program:

    ```bash
    go run main.go
    ```

3. The program will:
    - Create a `generated-qrs` directory
    - Generate QR codes for each entry
    - Create a PDF file (`generated-qrs/qr_codes.pdf`) containing all QR codes
    - Clean up temporary PNG files after PDF generation

## Output

-   The generated PDF will have QR codes arranged in a 3-column layout
-   Each QR code will be 50mm x 50mm
-   The unique ID will be displayed below each QR code
-   New pages will be automatically created when needed
