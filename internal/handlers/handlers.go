package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/render"
	"github.com/jung-kurt/gofpdf"
)

type Handler struct {
	DB *sql.DB
	// add pdf generator
	// status checker
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		DB: db,
		// add pdf generator
		// status checker
	}
}

type submitLinksReq struct {
	Links []string `json:"links"`
}

type submitLinksResp struct {
	Links    map[string]string `json:"links"`
	LinksNum int64             `json:"links_num"`
}

func (h *Handler) SubmitLinks(w http.ResponseWriter, r *http.Request) {
	var req submitLinksReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid json"})
		return
	}
	if len(req.Links) == 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "no links provided"})
		return
	}

	batchID, err := createIdAndLinks(h.DB, req.Links)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "internal"})
		return
	}

	respMap := make(map[string]string, len(req.Links))
	for _, l := range req.Links {
		respMap[l] = "queued"
	}

	resp := submitLinksResp{
		Links:    respMap,
		LinksNum: batchID,
	}
	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, resp)
}

func (h *Handler) GetReport(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("links_list")
	if q == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "links_list required"})
		return
	}
	parts := strings.Split(q, ",")
	var idNum []int64
	for _, p := range parts {
		if p == "" {
			continue
		}
		v, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "invalid links_list"})
			return
		}
		idNum = append(idNum, v)
	}

	linksInfo, err := gatherLinksForId(h.DB, idNum)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "internal"})
		return
	}

	pdfBytes, err := GeneratePDF(linksInfo)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "pdf generation failed"})
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="report.pdf"`)
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBytes)
}

func createIdAndLinks(db *sql.DB, links []string) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	res, err := tx.Exec(`INSERT INTO id DEFAULT VALUES`)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	batchID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	stmt, err := tx.Prepare(`
        INSERT INTO links (id, url, status)
        VALUES (?, ?, 'queued')
    `)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	for _, link := range links {
		if _, err := stmt.Exec(batchID, link); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	return batchID, tx.Commit()
}

func gatherLinksForId(db *sql.DB, id []int64) (map[string]string, error) {
	if len(id) == 0 {
		return map[string]string{}, nil
	}

	placeholders := make([]string, len(id))
	args := make([]interface{}, len(id))

	for i, num := range id {
		placeholders[i] = "?"
		args[i] = num
	}

	query := `SELECT url, status
		FROM links WHERE id IN (` + strings.Join(placeholders, ",") + `)
	`
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)

	for rows.Next() {
		var url, status string

		if err := rows.Scan(&url, &status); err != nil {
			return nil, err
		}
		result[url] = status
	}

	return result, nil
}

func GeneratePDF(data map[string]string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	pdf.Cell(40, 10, "Links Report")
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)

	for url, status := range data {
		line := fmt.Sprintf("%s — %s", url, status)
		pdf.Cell(0, 8, line)
		pdf.Ln(8)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
