package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/gen"
	"github.com/hubkudev/sentinel/internal/repositories"
	"github.com/xuri/excelize/v2"
)

type DownloadService interface {
	DownloadEventData(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, interval int32) ([]string, [][]string, error)
	WriteCSV(wFile *os.File, header []string, body [][]string) error
	WriteXLSX(wFile *os.File, header []string, body [][]string) error
	WriteJSON(wFile *os.File, header []string, body [][]string) error
	WriteHTML(wFile *os.File, header []string, body [][]string) error
	WritePDF(wFile *os.File, header []string, body [][]string) error
}

type DownloadServiceImpl struct {
	UtilService UtilService
	Repo        repositories.DownloadRepo
}

func InitDownloadService(utilService UtilService, repository repositories.DownloadRepo) DownloadServiceImpl {
	return DownloadServiceImpl{
		UtilService: utilService,
		Repo:        repository,
	}
}

func (s *DownloadServiceImpl) DownloadEventData(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, interval int32) ([]string, [][]string, error) {
	header, err := s.Repo.GetEventTableHeaders(ctx)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	input := gen.DownloadIntervalEventDataParams{
		UserID:    userID,
		ProjectID: projectID,
		Interval:  interval,
	}

	body, err := s.Repo.DownloadIntervalData(ctx, &input)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	var result [][]string
	for _, row := range body {
		var ip string

		if row.IpAddr != nil {
			ip = row.IpAddr.String()
		} else {
			ip = ""
		}

		item := []string{
			row.ID.String(),
			row.EventType,
			row.EventLabel.String,
			row.PageUrl.String,
			row.ElementPath.String,
			row.ElementType.String,
			ip,
			row.UserAgent.String,
			row.BrowserName.String,
			row.Country.String,
			row.Region.String,
			row.City.String,
			row.SessionID.String,
			row.DeviceType.String,
			fmt.Sprintf("%d", row.TimeOnPage.Int32),
			row.ScreenResolution.String,
			row.FiredAt.String(),
			row.ReceivedAt.String(),
			row.UserID.String(),
			row.ProjectID.String(),
		}

		result = append(result, item)
	}

	return header, result, nil
}

func (s *DownloadServiceImpl) WriteCSV(wFile *os.File, header []string, body [][]string) error {
	writer := csv.NewWriter(wFile)
	defer writer.Flush()

	if err := writer.Write(header); err != nil {
		return err
	}

	for _, row := range body {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	if err := writer.Error(); err != nil {
		return err
	}

	return nil
}

func (s *DownloadServiceImpl) WriteXLSX(wFile *os.File, header []string, body [][]string) error {
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	sheetName := "main"
	index, err := file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// set the header
	for idx, head := range header {
		cell := s.UtilService.LookupAlphabetFromIdx(idx)
		if err := file.SetCellValue(sheetName, fmt.Sprintf("%s1", cell), head); err != nil {
			return err
		}
	}

	// set the body
	for i, rows := range body {
		for j, row := range rows {
			cell := s.UtilService.LookupAlphabetFromIdx(j)

			// why + 2? because we want to start at row 2
			if err := file.SetCellValue(sheetName, fmt.Sprintf("%s%d", cell, i+2), row); err != nil {
				return err
			}
		}
	}

	file.SetActiveSheet(index)

	if err := file.Write(wFile); err != nil {
		return err
	}

	return nil
}

func (s *DownloadServiceImpl) WriteJSON(wFile *os.File, header []string, body [][]string) error {
	// preallocate json map based on header length
	jsonMap := make(map[string][]interface{}, len(header))

	// init header slices with capacity for body length
	for _, head := range header {
		jsonMap[head] = make([]interface{}, 0, len(body))
	}

	// fill the json body
	for _, rows := range body {
		for j, row := range rows {
			jsonMap[header[j]] = append(jsonMap[header[j]], row)
		}
	}

	// convert map to json
	result, err := json.Marshal(jsonMap)
	if err != nil {
		return err
	}

	// write to the file
	if _, err := wFile.Write(result); err != nil {
		return err
	}

	return nil
}

func (s *DownloadServiceImpl) WriteHTML(wFile *os.File, header []string, body [][]string) error {
	htmlBuilder := s.WriteHTMLBuilder(header, body)

	if _, err := wFile.WriteString(htmlBuilder); err != nil {
		return err
	}

	return nil
}

func (s *DownloadServiceImpl) WriteHTMLBuilder(header []string, body [][]string) string {
	var htmlBuilder strings.Builder

	// Write the opening tags for HTML and table
	htmlBuilder.WriteString(`
        <!DOCTYPE html>
        <html>
        <head>
            <title>Event Data</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    margin: 0px;
                }
                table {
                    width: 100%;
                    border-collapse: collapse;
                    margin-bottom: 0px;
                }
                th {
                    background-color: #f2f2f2;
                    padding: 2px;
                    text-align: left;
                    border: 1px solid #ddd;
                    font-weight: bold;
                }
                td {
                    padding: 2px;
                    border: 1px solid #ddd;
					max-width: 150px;
                    word-wrap: break-word;
                }
                tr:nth-child(even) {
                    background-color: #f9f9f9;
                }
                tr:hover {
                    background-color: #f5f5f5;
                }
            </style>
        </head>
        <body>
    `)
	htmlBuilder.WriteString("<table border='1'>\n")

	// Add table headers
	htmlBuilder.WriteString("<thead>\n<tr>\n")
	for _, head := range header {
		htmlBuilder.WriteString(fmt.Sprintf("<th>%s</th>\n", head))
	}
	htmlBuilder.WriteString("</tr>\n</thead>\n")

	// Add table rows
	htmlBuilder.WriteString("<tbody>\n")
	for _, rows := range body {
		htmlBuilder.WriteString("<tr>\n")
		for _, cell := range rows {
			htmlBuilder.WriteString(fmt.Sprintf("<td>%v</td>\n", cell))
		}
		htmlBuilder.WriteString("</tr>\n")
	}
	htmlBuilder.WriteString("</tbody>\n")

	// Close table and HTML tags
	htmlBuilder.WriteString("</table>\n</body>\n</html>")

	return htmlBuilder.String()
}

func (s *DownloadServiceImpl) WritePDF(wFile *os.File, header []string, body [][]string) error {
	htmlBuilder := s.WriteHTMLBuilder(header, body)

	// Create a new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return err
	}

	// add a new page from htmlBuilder
	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader([]byte(htmlBuilder))))

	// page settings
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeLegal)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationLandscape)
	pdfg.Dpi.Set(300)
	pdfg.MarginTop.Set(10)
	pdfg.MarginBottom.Set(10)

	// Create the PDF file
	if err := pdfg.Create(); err != nil {
		return err
	}

	// Write to file
	if _, err := wFile.Write(pdfg.Bytes()); err != nil {
		return fmt.Errorf("failed to write PDF to file: %w", err)
	}

	return nil
}
