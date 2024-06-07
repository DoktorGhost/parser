package csvRecord

import (
	"encoding/csv"
	"log"
	"os"
	"parser/internal/storage"
)

type CsvRecord struct{}

func NewCsvRecord() *CsvRecord {
	return &CsvRecord{}
}

func (c *CsvRecord) Add(name string, card storage.Card) error {
	// Open the CSV file in append mode
	file, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	// Initialize a file writer
	writer := csv.NewWriter(file)
	defer writer.Flush() // Ensure all data is written to the file

	// Check if the file is new (empty)
	fileInfo, err := os.Stat(name)
	if err != nil {
		log.Println("Error stating file:", err)
		return err
	}

	// Define the CSV headers, write them if the file is new
	if fileInfo.Size() == 0 {
		headers := []string{
			"Address",
			"Category",
			"SubCategory",
			"Name",
			"Url",
			"UrlImage",
			"Price",
			"PriceWithoutDiscount",
		}
		if err := writer.Write(headers); err != nil {
			log.Println("Error writing headers:", err)
			return err
		}
	}

	// Prepare the record to be written
	record := []string{
		card.Address,
		card.Category,
		card.SubCategory,
		card.Name,
		card.Url,
		card.UrlImage,
		card.Price,
		card.PriceWithoutDiscount,
	}

	// Write the record
	if err := writer.Write(record); err != nil {
		log.Println("Error writing record:", err)
		return err
	}

	return writer.Error()
}
