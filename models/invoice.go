package models

type Invoice struct {
	AccessNumber    string  `json:"custbody_access_number"`
	Serial          string  `json:"custbody_diag_invoice_serial"`
	Form            string  `json:"custbody_diag_invoice_form"`
	InvoiceNo       uint32  `json:"custbody_diag_invoice_no"`
	IssueDate       string  `json:"custbody_diag_invoice_date"`
	InvoiceAmt      float64 `json:"custbody_4"`
	ReservationCode string  `json:"custbody_7"`
	AdjustmentType  *string `json:"custbody_diag_invoice_adjtype"`
}

type AdjustInvoice struct {
	AccessNumber string `json:"custbody_access_number"`
	Serial       string `json:"custbody_diag_invoice_serial"`
	Form         string `json:"custbody_diag_invoice_form"`
	AdjustNo     uint32 `json:"custbody_9"`
}
