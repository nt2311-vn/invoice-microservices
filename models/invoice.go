package models

type Invoice struct {
	AccessNumber       string   `json:"custbody_access_number"`
	Serial             string   `json:"custbody_diag_invoice_serial"`
	Form               string   `json:"custbody_diag_invoice_form"`
	InvoiceNo          *uint32  `json:"custbody_diag_invoice_no"`
	IssueDate          *string  `json:"custbody_diag_invoice_date"`
	InvoiceAmt         *float64 `json:"custbody_4"`
	ReservationCode    *string  `json:"custbody_7"`
	AdjustNo           *uint32  `json:"custbody_9"`
	AdjustDate         *string  `json:"custbody_diag_invoice_adj_date"`
	AdjInvoiceAmt      *float64 `json:"custbody_diag_invadj_amt"`
	AdjustmentType     *string  `json:"custbody_diag_invoice_adjtype"`
	AdjReservationCode *string  `json:"custbody_vat_lookup_code_adj"`
}
